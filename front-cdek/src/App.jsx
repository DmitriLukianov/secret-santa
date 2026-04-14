import { BrowserRouter as Router, Routes, Route, Navigate, useParams } from 'react-router-dom'
import Home from './pages/home'
import Registration from './pages/registration'
import Registration_end from './pages/registration-end'
import Profile from './pages/profile'
import Profile_red from './pages/profile-red'
import Wishlist from './pages/wishlist'
import Wishlist_add from './pages/wishlist-add'
import Wishlist_red from './pages/wishlist-red'
import Wishlist_Santa from './pages/wishlist-santa'
import Game from './pages/game'
import Game_edit from './pages/game-edit'
import Game_add from './pages/game-add'
import Game_add_link from './pages/game-add-link'
import Letter from './pages/letter'
import SecretChat from './pages/chat'
import AuthCallback from './pages/auth-callback'
import InvitePage from './pages/invite'

function LegacyGameChatRedirect() {
  const { eventId } = useParams()
  return <Navigate to={`/game/${eventId}/chat`} replace />
}

function LegacyWishlistSantaNameRedirect() {
  const { name } = useParams()
  return <Navigate to={`/game/demo/wishlist/santa/${encodeURIComponent(name)}`} replace />
}

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/registration" element={<Registration />} />
        <Route path="/auth/callback" element={<AuthCallback />} />
        <Route path="/invite/:token" element={<InvitePage />} />
        <Route path="/join/:token" element={<InvitePage />} />
        <Route path="/registration-end" element={<Registration_end />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/profile/edit" element={<Profile_red />} />

        {/* Создание игры (без eventId) — статичные сегменты выше :eventId */}
        <Route path="/game/add/link" element={<Game_add_link />} />
        <Route path="/game/add" element={<Game_add />} />

        {/* Контекст конкретной игры */}
        <Route path="/game/:eventId/chat" element={<SecretChat />} />
        <Route path="/game/:eventId/edit" element={<Game_edit />} />
        <Route path="/game/:eventId/letter" element={<Letter />} />
        <Route path="/game/:eventId/wishlist/add" element={<Wishlist_add />} />
        <Route path="/game/:eventId/wishlist/items/:itemId" element={<Wishlist_red />} />
        <Route path="/game/:eventId/wishlist/santa/:participantSlug" element={<Wishlist_Santa />} />
        <Route path="/game/:eventId/wishlist/santa" element={<Wishlist_Santa />} />
        <Route path="/game/:eventId/wishlist" element={<Wishlist />} />
        <Route path="/game/:eventId" element={<Game />} />

        <Route path="/game" element={<Navigate to="/game/demo" replace />} />

        {/* Старые URL */}
        <Route path="/game-chat/:eventId" element={<LegacyGameChatRedirect />} />
        <Route path="/game-chat" element={<Navigate to="/game/demo/chat" replace />} />
        <Route path="/wishlist" element={<Wishlist />} />
        <Route path="/wishlist/add" element={<Wishlist_add />} />
        <Route path="/wishlist/items/:itemId" element={<Wishlist_red />} />
        <Route path="/wishlist-add" element={<Navigate to="/wishlist/add" replace />} />
        <Route path="/wishlist-red" element={<Navigate to="/game/demo/wishlist/items/1" replace />} />
        <Route path="/wishlist-santa/:name" element={<LegacyWishlistSantaNameRedirect />} />
        <Route path="/wishlist-santa" element={<Navigate to="/game/demo/wishlist/santa" replace />} />
        <Route path="/game-edit" element={<Navigate to="/game/demo/edit" replace />} />
        <Route path="/game-letter" element={<Navigate to="/game/demo/letter" replace />} />
        <Route path="/game-add-link" element={<Navigate to="/game/add/link" replace />} />
        <Route path="/game-add" element={<Navigate to="/game/add" replace />} />
        <Route path="/profile-red" element={<Navigate to="/profile/edit" replace />} />
      </Routes>
    </Router>
  )
}

export default App
