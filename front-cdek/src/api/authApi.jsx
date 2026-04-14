import { BASE_URL, getHeaders, handleResponse } from './http.jsx';

export const fetchMe = async () => {
  const response = await fetch(`${BASE_URL}/users/me`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const updateMe = async (userData) => {
  const response = await fetch(`${BASE_URL}/users/me`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify(userData),
  });
  return handleResponse(response);
};

export const sendOtpCode = async (email) => {
  const response = await fetch(`${BASE_URL}/auth/send-otp`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      // Токен здесь обычно не нужен, так как пользователь еще не вошел
    },
    body: JSON.stringify({ email }),
  });
  return handleResponse(response);
};


export const verifyOtpCode = async (name, email, code) => {
  const response = await fetch(`${BASE_URL}/auth/verify-otp`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ 
      name,
      email, 
      code 
    }),
  });
  
  const data = await handleResponse(response);
  
  // Если бэк возвращает токен сразу после верификации, сохраняем его
  if (data.token) {
    localStorage.setItem('token', data.token);
  }
  if (data.user) {
    localStorage.setItem('user', JSON.stringify(data.user));
  }
  
  return data;
};

/**
 * 3. Вход через GitHub (OAuth)
 * POST /api/v1/auth/login
 * Обычно OAuth работает через редирект, но если бэк ждет POST запрос с кодом авторизации:
 * @param {string} githubCode - Код, полученный от GitHub после редиректа
 */
export const loginWithGithub = async (githubCode) => {
  const response = await fetch(`${BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      provider: 'github', // Или как требует твой бэк
      code: githubCode
    }),
  });

  const data = await handleResponse(response);

  // Сохраняем токен и данные пользователя
  if (data.token) {
    localStorage.setItem('token', data.token);
  }
  if (data.user) {
    localStorage.setItem('user', JSON.stringify(data.user));
  }

  return data;
};
export const isAuthenticated = () => {
  const token = localStorage.getItem('token');
  if (!token) return false;
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    if (payload.exp && payload.exp < Date.now() / 1000) {
      localStorage.removeItem('token');
      return false;
    }
    return true;
  } catch {
    return false;
  }
};

export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};
