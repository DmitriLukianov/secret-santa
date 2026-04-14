export const BASE_URL = '/api/v1';

export const getHeaders = () => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
  };
};

export const handleResponse = async (response) => {
  if (!response.ok) {
    if (response.status === 401) {
      localStorage.removeItem('token');
    }
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || errorData.error || `Ошибка HTTP: ${response.status}`);
  }
  if (response.status === 204) return null;
  const text = await response.text();
  if (!text) return null;
  return JSON.parse(text);
};
