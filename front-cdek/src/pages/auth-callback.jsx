import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';

function AuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get('token');
    if (token) {
      localStorage.setItem('token', token);
      const pendingInvite = sessionStorage.getItem('pendingInviteToken');
      if (pendingInvite) {
        sessionStorage.removeItem('pendingInviteToken');
        navigate(`/invite/${pendingInvite}`, { replace: true });
      } else {
        navigate('/profile', { replace: true });
      }
    } else {
      navigate('/registration', { replace: true });
    }
  }, [searchParams, navigate]);

  return <div style={{ textAlign: 'center', marginTop: '100px' }}>Авторизация...</div>;
}

export default AuthCallback;
