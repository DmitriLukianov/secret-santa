// src/api/gameApi.js

const BASE_URL = '/api/v1';

// Вспомогательная функция для получения заголовков с токеном
const getHeaders = () => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
  };
};

// Вспомогательная функция для обработки ответов
const handleResponse = async (response) => {
  if (!response.ok) {
    // Пытаемся получить текст ошибки от сервера
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.message || `Ошибка HTTP: ${response.status}`);
  }
  // Если ответ пустой (например, 204 No Content), возвращаем null
  if (response.status === 204) return null;
  return response.json();
};

/* =========================================
   1. ПОЛЬЗОВАТЕЛЬ (AUTH / PROFILE)
   ========================================= */

// Получить данные текущего пользователя (для проверки авторизации и профиля)
export const fetchMe = async () => {
  const response = await fetch(`${BASE_URL}/users/me`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Обновить имя или email
export const updateMe = async (userData) => {
  const response = await fetch(`${BASE_URL}/users/me`, {
    method: 'PATCH',
    headers: getHeaders(),
    body: JSON.stringify(userData),
  });
  return handleResponse(response);
};

// Проверка авторизации (просто проверяем наличие токена локально)
export const isAuthenticated = () => {
  return !!localStorage.getItem('token');
};

// Выход из системы
export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

/* =========================================
   2. ИГРЫ (EVENTS) - Создание и Список
   ========================================= */

// Создать новую игру
export const createGame = async (gameData) => {
  const response = await fetch(`${BASE_URL}/events`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(gameData),
  });
  return handleResponse(response);
};

// Получить список игр текущего пользователя
export const fetchUserGames = async () => {
  const response = await fetch(`${BASE_URL}/events`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Получить конкретную игру по ID
export const fetchGameById = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Обновить игру (например, дату или название)
export const updateGame = async (eventId, updatedData) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(updatedData),
  });
  return handleResponse(response);
};

// Удалить игру
export const deleteGame = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

/* =========================================
   3. УЧАСТНИКИ (PARTICIPANTS)
   ========================================= */

// Добавить участника (обычно делает организатор)
export const addParticipant = async (eventId, participantData) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/participants`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(participantData),
  });
  return handleResponse(response);
};

// Получить список участников игры
export const fetchParticipants = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/participants`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Удалить участника
export const removeParticipant = async (participantId) => {
  const response = await fetch(`${BASE_URL}/participants/${participantId}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

/* =========================================
   4. ЖЕРЕБЬЁВКА И СТАТУСЫ
   ========================================= */

// Провести жеребьёвку
export const runDraw = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/assign`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Получить результаты жеребьёвки (кому я дарю)
export const fetchAssignments = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/assignments`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Активировать игру (начать)
export const activateGame = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/activate`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Завершить игру
export const finishGame = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/finish`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

/* =========================================
   5. ПРИГЛАШЕНИЯ И ПОДКЛЮЧЕНИЕ
   ========================================= */

// Сгенерировать ссылку-приглашение
export const generateInviteLink = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/open-invitation`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Вступить в игру по ссылке/коду
export const joinGameByLink = async (inviteCodeOrLink) => {
  // Предположим, что бэк ждет просто код или полную ссылку в поле invitationLink
  const response = await fetch(`${BASE_URL}/invite/join`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ invitationLink: inviteCodeOrLink }),
  });
  return handleResponse(response);
};

// Отправить приглашение на почту
export const sendInviteEmail = async (eventId, email) => {
  const response = await fetch(`${BASE_URL}/invitations/send-email`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ eventId, email }),
  });
  return handleResponse(response);
};

/* =========================================
   6. ВИШЛИСТЫ (WISHLISTS)
   ========================================= */

// Получить свой вишлист для конкретной игры
export const fetchMyWishlist = async (eventId) => {
  const response = await fetch(`${BASE_URL}/users/me/wishlist?eventId=${eventId}`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Получить вишлист другого участника (чтобы выбрать подарок)
export const fetchParticipantWishlist = async (participantId, eventId) => {
  const response = await fetch(`${BASE_URL}/wishlists/${participantId}?eventId=${eventId}`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Добавить товар в вишлист
export const addWishlistItem = async (wishlistId, itemData) => {
  const response = await fetch(`${BASE_URL}/wishlists/${wishlistId}/items`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(itemData),
  });
  return handleResponse(response);
};

// Удалить товар из вишлиста
export const deleteWishlistItem = async (wishlistId, itemId) => {
  const response = await fetch(`${BASE_URL}/wishlists/${wishlistId}/items/${itemId}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

/* =========================================
   7. ЧАТ (CHAT)
   ========================================= */

// Получить сообщения от того, кому я дарю
export const fetchRecipientChat = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/chat/recipient`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

// Отправить сообщение
export const sendMessage = async (eventId, messageText) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/chat/messages`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ text: messageText }),
  });
  return handleResponse(response);
};