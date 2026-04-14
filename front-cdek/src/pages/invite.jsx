import { useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { joinGameByLink, isAuthenticated } from '/src/api/gameApi.jsx';

function InvitePage() {
  const { token } = useParams();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated()) {
      sessionStorage.setItem('pendingInviteToken', token);
      navigate('/registration', { replace: true });
      return;
    }

    const join = async () => {
      try {
        const data = await joinGameByLink(token);
        navigate(`/game/${data.eventId}`, { replace: true });
      } catch {
        alert('Не удалось подключиться к игре. Попросите организатора прислать новую ссылку-приглашение.');
        navigate('/profile', { replace: true });
      }
    };

    join();
  }, [token, navigate]);

  return (
    <div style={{ textAlign: 'center', marginTop: '100px', color: '#757575' }}>
      Подключение к игре...
    </div>
  );
}

export default InvitePage;
