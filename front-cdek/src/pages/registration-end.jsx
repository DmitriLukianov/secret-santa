import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
// Импортируем функцию для получения данных пользователя
import { fetchMe, fetchUserGames } from '/src/api/gameApi.jsx';
import './main.css';

function Registration_end() {
  const navigate = useNavigate();
  
  // Состояния для данных
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState(null);
  const [games, setGames] = useState([]);

  // ← НОВОЕ: Проверка авторизации и загрузка данных при монтировании
  useEffect(() => {
    const checkAuthAndLoadData = async () => {
      try {
        setIsLoading(true);
        
        // 1. Проверяем, залогинен ли пользователь (получаем его данные)
        const userData = await fetchMe();
        setUser(userData);
        
        // 2. Получаем список игр пользователя, чтобы знать, куда вести на вишлист
        const gamesList = await fetchUserGames();
        const games = Array.isArray(gamesList) ? gamesList : (gamesList.items || []);
        setGames(games);
        
      } catch (error) {
        console.error('Ошибка загрузки данных:', error);
        // Если ошибка 401 (не авторизован), редиректим на регистрацию
        if (error.message.includes('401')) {
          navigate('/registration', { replace: true });
        }
      } finally {
        setIsLoading(false);
      }
    };

    checkAuthAndLoadData();
  }, [navigate]);

  const handleGoHome = () => {
    navigate('/');
  };

  const handleGoWishlist = () => {
    // Если есть игры, ведем на вишлист первой игры (или последней созданной)
    if (games.length > 0) {
      const firstGameId = games[0].id;
      navigate(`/game/${firstGameId}/wishlist`);
    } else {
      // Если игр нет, можно предложить создать игру или вести в профиль
      alert('У вас пока нет активных игр. Сначала создайте игру или подключитесь к существующей.');
      navigate('/profile');
    }
  };

  const handleGoProfile = () => {
    navigate('/profile'); 
  };

  // Рендер состояния загрузки
  if (isLoading) {
    return (
      <div className="overlay1">
        <div className="card1">
          <h1>Загрузка...</h1>
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <i className="ti ti-loader" style={{ fontSize: '48px', color: '#44E858', animation: 'spin 1s linear infinite' }}></i>
          </div>
        </div>
      </div>
    );
  }

  // Если пользователь не загрузился (ошибка авторизации), ничего не рендерим (редирект сработает в useEffect)
  if (!user) {
    return null;
  }

  return (
    <div className="overlay1">
      <div className="card1">
        <h1>ТАЙНЫЙ САНТА</h1>

        <form className="registration-end-form">
          {/* ← Динамическое приветствие */}
          <h3>Спасибо, {user.name || 'участник'}! Регистрация завершена.</h3>

          <button className="btn-primary" type="button" onClick={handleGoWishlist}>
            Создать вишлист
          </button>
          <button className="btn-primary" type="button" onClick={handleGoProfile}>
            Мой профиль
          </button>

          <button type="button" className="btn-secondary" onClick={handleGoHome}>
            На главную
          </button>
        </form>

      </div>
    </div>
  );
}

export default Registration_end;