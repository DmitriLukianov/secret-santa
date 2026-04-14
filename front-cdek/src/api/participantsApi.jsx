import { BASE_URL, getHeaders, handleResponse } from './http.jsx';

export const addParticipant = async (eventId, participantData) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/participants`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify(participantData),
  });
  return handleResponse(response);
};

export const fetchParticipants = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/participants`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const removeParticipant = async (participantId) => {
  const response = await fetch(`${BASE_URL}/participants/${participantId}`, {
    method: 'DELETE',
    headers: getHeaders(),
  });
  return handleResponse(response);
};
