import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { isAuthenticated } from '../api/authApi.jsx';
import { joinGameByLink } from '../api/invitationsApi.jsx';
import { fetchGameById, runDraw } from '../api/eventsApi.jsx';
import { fetchRecipientChat, sendMessage } from '../api/chatApi.jsx';
import './home.css';

void [joinGameByLink, fetchGameById, runDraw, fetchRecipientChat, sendMessage];

function Home() {
  // ... остальной код без изменений ...
  const [openIndex, setOpenIndex] = useState(null);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  
  const navigate = useNavigate();

  // Проверяем авторизацию при рендере
  const isAuth = isAuthenticated();

  const toggleFaq = (index) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  // Блокировка прокрутки при открытом меню
  useEffect(() => {
    document.body.style.overflow = isMobileMenuOpen ? 'hidden' : '';
    return () => { document.body.style.overflow = ''; };
  }, [isMobileMenuOpen]);

  // --- ЛОГИКА НАВИГАЦИИ С ПРОВЕРКОЙ АВТОРИЗАЦИИ ---

  // Кнопка "Профиль" (в шапке и моб. меню)
  const handleGoProfile = () => {
    if (isAuth) {
      navigate('/profile');
    } else {
      navigate('/registration'); 
    }
    setIsMobileMenuOpen(false);
  };

  // Кнопка "Вишлист"
  const handleGoWishlist = () => {
    if (isAuth) {
      navigate('/game/demo/wishlist');
    } else {
      navigate('/registration');
    }
    setIsMobileMenuOpen(false);
  };

  // Кнопка "НАЧАТЬ" (Hero и Rules)
  const handleStartGame = () => {
    if (isAuth) {
      // Если авторизован -> Сразу создаем игру
      navigate('/game/add');
    } else {
      // Если нет -> На регистрацию
      navigate('/registration');
    }
    setIsMobileMenuOpen(false);
  };

  const scrollToSection = (id) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
    setIsMobileMenuOpen(false);
  };

  const handleMobileNavClick = (action) => {
    setIsMobileMenuOpen(false);
    action();
  };

  const faqData = [
    { question: "Как долго длится игра?", answer: "Стандартно: 2–3 недели (от заполнения вишлистов до доставки). Даты задаёт организатор при создании." },
    { question: "Что, если участник не заполнил вишлист?", answer: "Вишлист может оставаться пустым по завершению жеребьевки. В таком случае Тайный Санта увидит пустой список и сможет выбрать подарок самостоятельно." },
    { question: "Когда я узнаю, кому дарить подарок?", answer: "В день жеребьёвки (дата задаётся при создании игры). Вы получите уведомление на почту с именем получателя и его вишлистом." },
    { question: "Могу ли я вытянуть самого себя?", answer: "Нет. Система гарантирует, что вы никогда не будете дарить подарок себе." }
  ];

  const rulesData = [
    { title: "Шаг 1. Регистрация", desc: "Чтобы начать, каждому участнику необходимо зарегистрироваться на сайте и войти в свой личный профиль." },
    { title: "Шаг 2. Создание игры", desc: "Зайдите в свой профиль и нажмите «Создать игру». Заполните данные: название команды, дату жеребьёвки и бюджет подарка. Скопируйте уникальную ссылку-приглашение или отправьте её прямо на email друзей через форму на сайте. Отправьте эту ссылку всем, кого хотите видеть в игре. После сохранения вы станете Организатором. Только вы сможете менять настройки игры, добавлять или удалять участников." },
    { title: "Шаг 3. Подключение к игре (для Участников)", desc: "Получив ссылку от друга, убедитесь, что вы зарегистрированы и авторизованы. Зайдите в свой профиль и нажмите кнопку «Подключиться к игре». Вставьте полученную ссылку" },
    { title: "Шаг 4. Вишлисты и Жеребьёвка", desc: "Каждый участник заполняет свой Вишлист (список желаемых подарков) в профиле. Это поможет вашему Тайному Санте выбрать идеальный презент. В назначенную организатором дату система автоматически проведёт Жеребьёвку и распределит пары. Вам будут доступны результаты жеребьевки: имя и вишлист человека, кому вы дарите подарок." },
    { title: "Шаг 5. Подготовьте подарок и вручите в праздничный день! ", desc: "" }
  ];

  const advantagesData = [
    { icon: "ti ti-truck-delivery", title: "Доставка через СДЭК", desc: "Доставим ваш подарок курьером дедом-морозом" },
    { icon: "ti ti-heart", title: "Вишлисты", desc: "Добавляйте ссылки на желаемые товары" },
    { icon: "ti ti-message-question", title: "Анонимные вопросы", desc: "Задавайте вопросы адресату в секретном чате" }
  ];

  return (
    <div className="app">
      {/*  HEADER  */}
        <header className="header">
          <nav>
            {/* Бургер */}
            <button 
              className={`burger-menu ${isMobileMenuOpen ? 'active' : ''}`}
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              aria-label="Меню"
            >
              <span></span><span></span><span></span>
            </button>
            <div className="header-left">
              <a href="#" onClick={(e) => { e.preventDefault(); handleGoProfile(); }}>Профиль</a>
              <a href="#" onClick={(e) => { e.preventDefault(); handleGoWishlist(); }}>Вишлист</a>
            </div>

            <a href="#" className="logo" onClick={(e) => { 
              e.preventDefault(); 
              window.scrollTo({ top: 0, behavior: 'smooth' }); 
            }}>
              <img src="/sdek-logo.png" alt="СДЭК" className="logo-img"
                  onError={(e) => {
                    e.currentTarget.style.display = 'none';
                    const parent = e.currentTarget.closest('.logo, .footer-logo');
                    const fallback = parent?.querySelector('.logo-fallback');
                    if (fallback) fallback.style.display = 'inline-block';
                  }} 
              />
              <span className="logo-fallback">СДЭК</span>
            </a>

            <div className="header-right">
              <a href="#rules" onClick={() => scrollToSection('rules')}>Правила</a>
              <a href="#faq" onClick={() => scrollToSection('faq')}>Вопросы</a>
            </div>
          </nav>
</header>

      {/*  МОБИЛЬНОЕ МЕНЮ  */}
        {isMobileMenuOpen && (
          <div className="mobile-overlay" onClick={() => setIsMobileMenuOpen(false)} />
        )}

        <div className={`mobile-nav ${isMobileMenuOpen ? 'active' : ''}`}>
      
          <a href="#" onClick={(e) => { 
            e.preventDefault(); 
            handleMobileNavClick(handleGoProfile); 
          }}>Профиль</a>
          
          <a href="#" onClick={(e) => { 
            e.preventDefault(); 
            handleMobileNavClick(handleGoWishlist); 
          }}>Вишлист</a>
          

          <a href="#" onClick={(e) => { 
            e.preventDefault(); 
            handleMobileNavClick(() => scrollToSection('rules')); 
          }}>Правила</a>
          
          <a href="#" onClick={(e) => { 
            e.preventDefault(); 
            handleMobileNavClick(() => scrollToSection('faq')); 
          }}>Вопросы</a>
        </div>

      {/*  HERO  */}
      <section className="hero">
        <h1>ТАЙНЫЙ САНТА</h1>
        <p>
          Один клик — и вы в игре. Один конверт — и вы узнаёте,<br />
          чью жизнь сделаете чуточку ярче этим Новым годом!
        </p>
        {/* ИЗМЕНЕНО: вызов handleStartGame вместо handleGoRegistration */}
        <button onClick={handleStartGame}>НАЧАТЬ</button>
      </section>

      {/*  ADVANTAGES  */}
      <section className="advantages">
        <div className="cards">
          {advantagesData.map((item, index) => (
            <div className="card" key={index}>
              <i className={item.icon} style={{ fontSize: '48px', color: '#44E858', marginBottom: '20px' }}></i>
              <h3>{item.title}</h3>
              <p>{item.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/*  RULES  */}
      <section className="rules" id="rules">
        <img 
          src="/cookie.png" 
          alt="Печенье" 
          className="decoration cookie"
        />
        <h2>ПРАВИЛА ИГРЫ</h2>
        <ul>
          {rulesData.map((rule, index) => (
            <li key={index}>
              <strong>{rule.title}</strong>
              {rule.desc && <span>{rule.desc}</span>}
            </li>
          ))}
        </ul>
        {/* ИЗМЕНЕНО: вызов handleStartGame вместо handleGoRegistration */}
        <button onClick={handleStartGame}>НАЧАТЬ</button>
      </section>

      <img 
        src="/garland.png" 
        alt="Гирлянда" 
        className="decoration garland"
      />

      {/*  FAQ  */}
      <section className="faq" id="faq">
        <h2>ЧАСТЫЕ ВОПРОСЫ</h2>
        <div className="faq-list">
          {faqData.map((item, index) => (
            <details key={index} className="faq-item" open={openIndex === index}>
              <summary onClick={(e) => { e.preventDefault(); toggleFaq(index); }}>
                {item.question}
              </summary>
              <p>{item.answer}</p>
            </details>
          ))}
        </div>
      </section>

      {/*  FOOTER  */}
      <footer className="footer">

        <img 
          src="/santa.png" 
          alt="Санта" 
          className="decoration santa"
        />
        <div className="footer-content">
          <div className="footer-section">
            <h4>Контакты</h4>
            <p>+7 (495) 009 04 05</p>
          </div>
          <div className="footer-section">
            <h4>Поддержка</h4>
            <a href="#faq">Частые вопросы</a>
          </div>
          <div className="footer-section">
            <h4>Правила</h4>
            <a href="#rules">Правила игры</a>
          </div>
          <div className="footer-logo">
            <img 
              src="/sdek-logo.png" 
              alt="СДЭК" 
              className="footer-logo-img"
              onError={(e) => {
                e.currentTarget.style.display = 'none';
                const parent = e.currentTarget.closest('.logo, .footer-logo');
                const fallback = parent?.querySelector('.logo-fallback');
                if (fallback) fallback.style.display = 'inline-block';
              }} 
            />
            <span className="logo-fallback">СДЭК</span>
          </div>
        </div>
        <div className="footer-bottom">© 2026 Тайный Санта. Все права защищены.</div>
      </footer>
      
    </div>
  );
}

export default Home;