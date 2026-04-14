import { BASE_URL, getHeaders, handleResponse } from './http.jsx';

export const fetchRecipientChat = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/chat/recipient`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const fetchSenderChat = async (eventId) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/chat/sender`, {
    method: 'GET',
    headers: getHeaders(),
  });
  return handleResponse(response);
};

export const sendMessage = async (eventId, messageText) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/chat/messages`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ text: messageText }),
  });
  return handleResponse(response);
};

export const sendMessageToSanta = async (eventId, messageText) => {
  const response = await fetch(`${BASE_URL}/events/${eventId}/chat/messages/santa`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ text: messageText }),
  });
  return handleResponse(response);
};
