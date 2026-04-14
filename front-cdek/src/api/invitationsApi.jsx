import { BASE_URL, getHeaders, handleResponse } from './http.jsx';

export const generateInviteLink = async (eventId) => {
  const response = await fetch(`${BASE_URL}/invitations/generate`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ eventId }),
  });
  return handleResponse(response);
};

export const joinGameByLink = async (inviteCodeOrLink) => {
  // Принимаем и полный URL, и просто токен
  let token;
  if (inviteCodeOrLink.includes('/invite/')) {
    token = inviteCodeOrLink.split('/invite/').pop().split('?')[0];
  } else if (inviteCodeOrLink.includes('/join/')) {
    token = inviteCodeOrLink.split('/join/').pop().split('?')[0];
  } else {
    token = inviteCodeOrLink.trim();
  }

  const response = await fetch(`${BASE_URL}/invite/join`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ invitationLink: token }),
  });
  return handleResponse(response);
};

export const sendInviteEmail = async (eventId, email) => {
  const response = await fetch(`${BASE_URL}/invitations/send-email`, {
    method: 'POST',
    headers: getHeaders(),
    body: JSON.stringify({ eventId, email }),
  });
  return handleResponse(response);
};
