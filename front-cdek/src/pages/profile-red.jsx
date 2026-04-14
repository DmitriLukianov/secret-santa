import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
// Импортируем методы API
import { fetchMe, updateMe, isAuthenticated } from '/src/api/gameApi.jsx';
import './main.css';

// === ФУНКЦИИ ВАЛИДАЦИИ (без изменений) ===
const validateName = (name) => {
  const errors = [];
  const trimmed = name.trim();
  if (!trimmed) {
    errors.push('Имя не может быть пустым');
    return errors;
  }
  const validPattern = /^[а-яА-ЯёЁa-zA-Z\s-]+$/;
  if (!validPattern.test(trimmed)) {
    errors.push('Имя может содержать только буквы, пробелы и дефисы');
  }
  if (trimmed.length < 2) {
    errors.push('Минимальная длина имени — 2 символа');
  }
  if (trimmed.length > 50) {
    errors.push('Максимальная длина имени — 50 символов');
  }
  if (trimmed.startsWith(' ') || trimmed.endsWith(' ') || 
      trimmed.startsWith('-') || trimmed.endsWith('-')) {
    errors.push('Имя не должно начинаться или заканчиваться пробелом или дефисом');
  }
  if (/\s{2,}/.test(trimmed)) {
    errors.push('Не должно быть нескольких пробелов подряд');
  }
  if (/--+/.test(trimmed)) {
    errors.push('Не должно быть нескольких дефисов подряд');
  }
  return errors;
};

const validateEmail = (email) => {
  const errors = [];
  if (!email) {
    errors.push('Почта обязательна для заполнения');
    return errors;
  }
  const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  if (!emailPattern.test(email)) {
    errors.push('Введите корректный адрес почты (пример: name@example.com)');
  }
  return errors;
};

function Profile_red() {
  const navigate = useNavigate();
  
  // Состояния для данных формы
  const [formData, setFormData] = useState({
    name: '',
    mail: ''
  });

  // Состояния загрузки и ошибок
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [submitError, setSubmitError] = useState(null);
  const [fetchError, setFetchError] = useState(null);

  // Состояния валидации
  const [errors, setErrors] = useState({ name: [], mail: [] });
  const [touched, setTouched] = useState({ name: false, mail: false });

  // ← НОВОЕ: Загрузка данных профиля при монтировании
  useEffect(() => {
    if (!isAuthenticated()) {
      navigate('/registration', { replace: true });
      return;
    }

    const loadProfile = async () => {
      try {
        setIsLoading(true);
        setFetchError(null);
        
        const user = await fetchMe();
        
        // Заполняем форму данными с сервера
        // Адаптируйте поля (user.name, user.email) под ваш ответ API
        setFormData({
          name: user.name || '',
          mail: user.email || ''
        });
      } catch (err) {
        console.error('Ошибка загрузки профиля:', err);
        setFetchError(err.message || 'Не удалось загрузить данные профиля');
        // Если ошибка 401, можно редиректить на вход
        if (err.message.includes('401')) {
          navigate('/', { replace: true });
        }
      } finally {
        setIsLoading(false);
      }
    };

    loadProfile();
  }, [navigate]);

  // Обработчик изменения полей
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
    setSubmitError(null); // Сбрасываем ошибку сервера при редактировании
    
    if (touched[name]) {
      const validator = name === 'name' ? validateName : validateEmail;
      setErrors({ ...errors, [name]: validator(value) });
    }
  };

  // Обработчик потери фокуса
  const handleBlur = (e) => {
    const { name, value } = e.target;
    setTouched({ ...touched, [name]: true });
    
    const validator = name === 'name' ? validateName : validateEmail;
    setErrors({ ...errors, [name]: validator(value) });
  };

  // Проверка валидности формы
  const isFormValid = () => {
    const nameErrors = validateName(formData.name);
    const mailErrors = validateEmail(formData.mail);
    setErrors({ name: nameErrors, mail: mailErrors });
    return nameErrors.length === 0 && mailErrors.length === 0;
  };

  const handleGoProfile = () => {
    navigate('/profile'); 
  };

  // ← НОВОЕ: Отправка данных на сервер
  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitError(null);

    if (!isFormValid()) {
      setTouched({ name: true, mail: true });
      return;
    }

    try {
      setIsSaving(true);

      // Подготовка данных для отправки
      // Бэкенд может ожидать { name, email } или { fullName, email } - адаптируйте под свой API
      const userData = {
        name: formData.name.trim(),
        email: formData.mail.trim()
      };

      // Вызов API обновления
      await updateMe(userData);

      alert('Данные успешно сохранены!');
      handleGoProfile();
      
    } catch (err) {
      console.error('Ошибка сохранения:', err);
      setSubmitError(err.message || 'Не удалось сохранить изменения. Попробуйте позже.');
      alert(err.message || 'Ошибка при сохранении');
    } finally {
      setIsSaving(false);
    }
  };
  
  const handleGoBack = () => {
    navigate(-1);
  };

  // Рендер состояния загрузки
  if (isLoading) {
    return (
      <div className="overlay_profile_red">
        <div className="card_profile_red">
          <h1>Загрузка...</h1>
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <i className="ti ti-loader" style={{ fontSize: '48px', color: '#44E858', animation: 'spin 1s linear infinite' }}></i>
          </div>
        </div>
      </div>
    );
  }

  // Рендер ошибки загрузки
  if (fetchError) {
    return (
      <div className="overlay_profile_red">
        <div className="card_profile_red">
          <h1>Ошибка</h1>
          <div style={{ textAlign: 'center', padding: '40px', color: '#e74c3c' }}>
            <p>{fetchError}</p>
            <button className="btn-secondary" onClick={handleGoBack} style={{ marginTop: '20px' }}>Назад</button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="overlay_profile_red"> 
      <div className="card_profile_red"> 
        <h1>Редактирование профиля</h1>

        <form onSubmit={handleSubmit} noValidate> 
          <div className="form-group">
            <label htmlFor="name">Имя *</label>
            <input 
              id="name"
              type="text" 
              name="name"
              placeholder="Введите имя" 
              value={formData.name}
              onChange={handleChange}
              onBlur={handleBlur}
              required 
              disabled={isSaving}
              className={`${errors.name.length > 0 && touched.name ? 'input-error' : ''} ${isSaving ? 'input-disabled' : ''}`}
              pattern="^[а-яА-ЯёЁa-zA-Z\s-]+$"
              minLength={2}
              maxLength={50}
            />
            {errors.name.length > 0 && touched.name && (
              <ul className="error-list">
                {errors.name.map((err, i) => (
                  <li key={i} className="error-item">• {err}</li>
                ))}
              </ul>
            )}
          </div>
          
          <div className="form-group">
            <label htmlFor="mail">Почта *</label>
            <input 
              id="mail"
              type="email" 
              name="mail"
              placeholder="Введите почту" 
              value={formData.mail}
              onChange={handleChange}
              onBlur={handleBlur}
              required
              disabled={isSaving}
              className={`${errors.mail.length > 0 && touched.mail ? 'input-error' : ''} ${isSaving ? 'input-disabled' : ''}`}
            />
            {errors.mail.length > 0 && touched.mail && (
              <ul className="error-list">
                {errors.mail.map((err, i) => (
                  <li key={i} className="error-item">• {err}</li>
                ))}
              </ul>
            )}
          </div>

          {submitError && (
            <div style={{ color: '#e74c3c', fontSize: '14px', marginBottom: '15px', textAlign: 'center' }}>
              ⚠️ {submitError}
            </div>
          )}

          <div className="button-group">
            <button 
              type="submit" 
              className="btn-primary" 
              disabled={isSaving}
              style={{ opacity: isSaving ? 0.7 : 1, cursor: isSaving ? 'not-allowed' : 'pointer' }}
            >
              {isSaving ? 'Сохранение...' : 'Сохранить'}
            </button>
            <button 
              type="button" 
              className="btn-secondary" 
              onClick={handleGoBack}
              disabled={isSaving}
            >
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Profile_red;