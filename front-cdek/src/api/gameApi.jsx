/**
 * Общая точка входа: реэкспорт всех методов API (обратная совместимость).
 * По доменам см. отдельные модули: authApi, eventsApi, participantsApi и т.д.
 */

export {
  fetchMe,
  updateMe,
  isAuthenticated,
  logout,
  sendOtpCode,
  verifyOtpCode,
  loginWithGithub,
} from './authApi.jsx';

export {
  createGame,
  fetchUserGames,
  fetchGameById,
  updateGame,
  deleteGame,
  runDraw,
  fetchAssignments,
  activateGame,
  finishGame,
} from './eventsApi.jsx';

export { addParticipant, fetchParticipants, removeParticipant } from './participantsApi.jsx';

export { generateInviteLink, joinGameByLink, sendInviteEmail } from './invitationsApi.jsx';

export {
  fetchMyWishlist,
  fetchWishlistItems,
  fetchParticipantWishlist,
  addWishlistItem,
  deleteWishlistItem,
  updateWishlistItem,
  uploadFile,
} from './wishlistApi.jsx';

export { fetchRecipientChat, fetchSenderChat, sendMessage, sendMessageToSanta } from './chatApi.jsx';

export { createProduct, fetchProducts } from './productsApi.jsx';
