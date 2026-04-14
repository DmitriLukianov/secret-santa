import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
// Импортируем нужные методы API
import { generateInviteLink, sendInviteEmail } from '/src/api/gameApi.jsx';
import './main.css';

function Game_add_link() {
    const navigate = useNavigate();
    const location = useLocation();
    // Получаем ID игры из state (передаётся при создании) или используем заглушку
    const createdEventId = location.state?.eventId;
    
    // Состояния модального окна
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isCopied, setIsCopied] = useState(false);
    
    // Состояния для данных
    const [inviteLink, setInviteLink] = useState('');
    const [isLoadingLink, setIsLoadingLink] = useState(false);
    
    // Состояние для email и ошибки
    const [email, setEmail] = useState('');
    const [emailError, setEmailError] = useState('');
    const [isSendingEmail, setIsSendingEmail] = useState(false);

    // ← НОВОЕ: Загрузка ссылки-приглашения при монтировании
    useEffect(() => {
        const loadInviteLink = async () => {
            if (!createdEventId) return;

            try {
                setIsLoadingLink(true);
                const data = await generateInviteLink(createdEventId);
                
                const url = data.inviteUrl || (data.token ? `${window.location.origin}/invite/${data.token}` : null);
                if (url) {
                    setInviteLink(url);
                }
            } catch (error) {
                console.error('Ошибка получения ссылки:', error);
            } finally {
                setIsLoadingLink(false);
            }
        };

        if (createdEventId) {
            loadInviteLink();
        }
    }, [createdEventId]);

    // Копировать ссылку
    const handleCopyLink = async () => {
        if (!inviteLink) return;
        try {
            await navigator.clipboard.writeText(inviteLink);
            setIsCopied(true);
            setTimeout(() => setIsCopied(false), 2000);
        } catch (err) {
            alert('Не удалось скопировать');
        }
    };

    // ← НОВОЕ: Проверка и отправка email через API
    const handleSendEmail = async () => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        
        if (!email) {
            setEmailError('Введите email');
            return;
        }
        
        if (!emailRegex.test(email)) {
            setEmailError('Некорректный email');
            return;
        }

        if (!createdEventId) {
            setEmailError('ID игры не найден');
            return;
        }

        try {
            setIsSendingEmail(true);
            setEmailError('');
            
            // Вызов API отправки приглашения
            await sendInviteEmail(createdEventId, email);
            
            alert(`Приглашение отправлено на ${email}`);
            setEmail('');
            handleCloseModal();
        } catch (error) {
            console.error('Ошибка отправки email:', error);
            setEmailError(error.message || 'Не удалось отправить приглашение');
        } finally {
            setIsSendingEmail(false);
        }
    };

    // Открыть окно
    const handleAddParticipant = () => {
        setIsModalOpen(true);
        setEmailError('');
        setEmail('');
    };

    // Закрыть окно
    const handleCloseModal = () => {
        setIsModalOpen(false);
        setEmailError('');
        setEmail('');
    };

    // Если ID игры нет, показываем ошибку или редиректим
    if (!createdEventId && !isLoadingLink) {
        return (
            <div className="overlay_profile_red">
                <div className="card_profile_red">
                    <h1>Ошибка</h1>
                    <p style={{ textAlign: 'center', color: '#757575' }}>Игра не найдена.</p>
                    <button className="btn-primary" onClick={() => navigate('/profile')}>В профиль</button>
                </div>
            </div>
        );
    }

    return (
        <div className="overlay_profile_red"> 
            <div className="card_profile_red"> 
                <h1>Создание игры</h1>
                
                {/* Обновленный параграф с стилями */}
                <p style={{ 
                    textAlign: 'center', 
                    color: '#757575', 
                    fontSize: '14px', 
                    lineHeight: '1.5',
                    marginBottom: '20px',
                    maxWidth: '350px',
                    margin: '0 auto 20px auto'
                }}>
                    Игра успешно создана! Пригласите участников по ссылке или отправьте приглашение на email.
                </p>

                {/* Кнопка открытия модалки */}
                <button 
                    type="button" 
                    className="btn-secondary"
                    onClick={handleAddParticipant}
                    disabled={!inviteLink}
                    style={{ width: '100%', maxWidth: '350px', margin: '0 auto 20px auto', display: 'block', opacity: !inviteLink ? 0.6 : 1 }}
                >
                    + Добавить участников
                </button>

                <div className="button-group" style={{ flexDirection: 'column' }}> 
                    <button type="button" className="btn-primary" onClick={() => navigate(`/game/${createdEventId}`)}>
                        Перейти к игре
                    </button>
                    <button type="button" className="btn-secondary" onClick={() => navigate('/profile')}>
                        Мой профиль
                    </button>
                </div>
            </div>

            {/* Модальное окно */}
            {isModalOpen && (
                <div className="modal-overlay" onClick={handleCloseModal}>
                    <div className="modal-small" onClick={(e) => e.stopPropagation()}>
                        <button className="modal-close" onClick={handleCloseModal}>×</button>
                        
                        <h3>Приглашение участников</h3>

                        {/* Блок 1: Ссылка */}
                        <p className="modal-label">Ссылка для приглашения:</p>
                        <div className="link-row">
                            <input 
                                type="text" 
                                className="link-input" 
                                value={isLoadingLink ? 'Загрузка...' : inviteLink} 
                                readOnly 
                                disabled={isLoadingLink}
                            />
                            <button 
                                type="button" 
                                className="btn-primary"
                                onClick={handleCopyLink}
                                disabled={!inviteLink || isLoadingLink || isCopied}
                                style={{ minWidth: '100px' }}
                            >
                                {isCopied ? '✓ Скопировано!' : 'Копировать'}
                            </button>
                        </div>

                        {/* Блок 2: Email с проверкой */}
                        <p className="modal-label" style={{ marginTop: '15px' }}>Отправить приглашение на почту:</p>
                        <div className="link-row">
                            <input 
                                type="email" 
                                className={`link-input ${emailError ? 'input-error' : ''}`}
                                placeholder="example@mail.com"
                                value={email}
                                onChange={(e) => {
                                    setEmail(e.target.value);
                                    if (emailError) setEmailError('');
                                }}
                                disabled={isSendingEmail}
                            />
                            <button 
                                type="button" 
                                className="btn-primary"
                                onClick={handleSendEmail}
                                disabled={isSendingEmail}
                                style={{ minWidth: '100px', opacity: isSendingEmail ? 0.7 : 1 }}
                            >
                                {isSendingEmail ? 'Отправка...' : 'Отправить'}
                            </button>
                        </div>
                        
                        {/* Текст ошибки */}
                        {emailError && <span className="error-text" style={{ color: '#e74c3c', fontSize: '12px', marginTop: '5px', display: 'block' }}>{emailError}</span>}
                        
                    </div>
                </div>
            )}
        </div>
    );
}

export default Game_add_link;