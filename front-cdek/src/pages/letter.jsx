import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { fetchAssignments, fetchParticipants, fetchMe, isAuthenticated } from '/src/api/gameApi.jsx';
import './main.css';

function Letter({ organizerMessage }) {
  const navigate = useNavigate();
  const { eventId } = useParams();

  const [myName, setMyName] = useState('');
  const [recipientName, setRecipientName] = useState('');
  const [recipientParticipantId, setRecipientParticipantId] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const handleGoWishlist = () => {
    if (recipientParticipantId) {
      navigate(`/game/${eventId}/wishlist/santa/${recipientParticipantId}?name=${encodeURIComponent(recipientName)}`);
    } else {
      navigate(`/game/${eventId}/wishlist/santa`);
    }
  };

  useEffect(() => {
    if (!isAuthenticated()) {
      navigate('/registration', { replace: true });
      return;
    }

    const loadData = async () => {
      if (!eventId) return;

      try {
        setIsLoading(true);

        // Параллельная загрузка: текущий пользователь, назначения, участники
        const [me, assignments, participants] = await Promise.all([
          fetchMe(),
          fetchAssignments(eventId),
          fetchParticipants(eventId),
        ]);

        setMyName(me?.name || '');

        // fetchAssignments возвращает массив; берём первый элемент
        const assignment = Array.isArray(assignments) ? assignments[0] : null;

        if (assignment) {
          setRecipientName(assignment.receiverName || '');

          // Найдём participantId получателя в списке участников
          const participantsList = Array.isArray(participants) ? participants : [];
          const recipientParticipant = participantsList.find(
            p => p.userId === assignment.receiverId
          );
          if (recipientParticipant) {
            setRecipientParticipantId(recipientParticipant.id);
          }
        }
      } catch (error) {
        console.error('Ошибка загрузки данных письма:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [eventId, navigate]);

  if (isLoading) {
    return (
      <div className="letter-page">
        <div className="letter-card">
          <button className="letter-close" onClick={() => navigate(-1)}>×</button>
          <p style={{ textAlign: 'center', paddingTop: '50px' }}>Загрузка письма...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="letter-page">
      <div className="letter-card">
        <button className="letter-close" onClick={() => navigate(-1)}>×</button>

        <div className="letter-header">
          <div className="letter-from">
            <span className="label">От:</span>
            <span className="name">{recipientName || 'Участник'}</span>
          </div>
          <div className="letter-to">
            <span className="label">Кому:</span>
            <span className="name">Санте ({myName || 'вам'})</span>
          </div>
        </div>

        <div className="letter-body">
          <p className="greeting">Дорогой Санта,</p>

          <p className="letter-text">
            Меня зовут <span className="highlight">{recipientName || 'Участник'}</span>.
            В этом году я стабильно показывал(а) хорошие результаты!
            В качестве подарка я буду рад(а) получить следующее:
          </p>

          <button className="btn-primary" onClick={handleGoWishlist}>
            Вишлист получателя
          </button>

          <p className="letter-text">
            Обещаю в новом году быть еще ответственнее, помогать коллегам и верить в чудо,
            даже когда дедлайны горят!
          </p>

          <p className="closing">
            С праздничным настроением, <span className="highlight">{recipientName || 'Участник'}</span>
          </p>

          {organizerMessage && (
            <div className="organizer-message">
              <div className="organizer-badge">
                <i className="ti ti-speakerphone" style={{ fontSize: '16px', marginRight: '6px' }}></i>
                От организатора
              </div>
              <p className="organizer-text">{organizerMessage}</p>
            </div>
          )}
        </div>

        <div className="letter-border"></div>
      </div>
    </div>
  );
}

export default Letter;
