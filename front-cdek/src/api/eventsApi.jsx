import { BASE_URL, getHeaders, handleResponse } from './http.jsx';

export const createGame = async (gameData) => {
  const response = await fetch(`${BASE_URL}/events`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(gameData),
  });
  return handleResponse(response);
};

export const fetchUserGames = async () => {
  const response = await fetch(`${BASE_URL}/events`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const fetchGameById = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const updateGame = async (eventId, updatedData) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}`, {
    method: 'PUT',
    headers: getHeaders(),
    body: JSON.stringify(updatedData),
  });
  return handleResponse(response);
};

export const deleteGame = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const runDraw = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/assign`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const fetchAssignments = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/assignments`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const activateGame = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/activate`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const finishGame = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/finish`, {
    method: 'POST',
    headers: getHeaders(),
  });
  return handleResponse(response);
};
