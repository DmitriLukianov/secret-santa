import React, { useState, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
// Импортируем нужные методы из вашего API
import { fetchMyWishlist, addWishlistItem, isAuthenticated, uploadFile } from '/src/api/gameApi.jsx';
import './main.css';

// === ФУНКЦИИ ВАЛИДАЦИИ (без изменений) ===
const validateName = (name) => {
  const errors = [];
  const trimmed = name.trim();
  if (!trimmed) {
    errors.push('Название товара обязательно');
    return errors;
  }
  const validPattern = /^[а-яА-ЯёЁa-zA-Z0-9\s\-\,\.\(\)\/]+$/;
  if (!validPattern.test(trimmed)) {
    errors.push('Разрешены только буквы, цифры, пробел и символы: - , . ( ) /');
  }
  if (trimmed.length < 3) {
    errors.push('Минимальная длина названия — 3 символа');
  }
  if (trimmed.length > 150) {
    errors.push('Максимальная длина названия — 150 символов');
  }
  if (trimmed.startsWith(' ') || trimmed.endsWith(' ')) {
    errors.push('Название не должно начинаться или заканчиваться пробелом');
  }
  return errors;
};

const validatePrice = (price) => {
  const errors = [];
  const num = parseFloat(price);
  if (isNaN(num)) {
    errors.push('Цена должна быть числом');
    return errors;
  }
  if (num <= 0) {
    errors.push('Цена должна быть больше 0');
  }
  if (num > 1000000) {
    errors.push('Максимальная цена — 1 000 000');
  }
  const parts = price.toString().split('.');
  if (parts[1] && parts[1].length > 2) {
    errors.push('Цена может иметь не более 2 знаков после запятой');
  }
  return errors;
};

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

function WishlistAdd() {
  const navigate = useNavigate();
  const { eventId: rawEventId } = useParams();
  const eventId = rawEventId && UUID_RE.test(rawEventId) ? rawEventId : undefined;

  React.useEffect(() => {
    if (!isAuthenticated()) {
      navigate('/registration', { replace: true });
    }
  }, [navigate]);
  const fileInputRef = useRef(null);

  const [formData, setFormData] = useState({
    name: '',
    price: '',
    link: ''
  });

  const [errors, setErrors] = useState({ name: [], price: [] });
  const [touched, setTouched] = useState({ name: false, price: false });
  
  // ← НОВОЕ: Состояния для загрузки и ошибок сервера
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState(null);

  const [isDragging, setIsDragging] = useState(false);
  const [files, setFiles] = useState([]);

  const handleGoWishlist = () => {
    if (eventId) {
      navigate(`/game/${eventId}/wishlist`);
    } else {
      navigate('/wishlist');
    }
  };

  const handleGoBack = () => {
    navigate(-1);
  };

  // Обработчики полей (без изменений)
  const handlePriceChange = (e) => {
    const value = e.target.value;
    if (value === '' || /^\d*\.?\d*$/.test(value)) {
      setFormData(prev => ({ ...prev, price: value }));
      if (touched.price) {
        setErrors(prev => ({ ...prev, price: validatePrice(value) }));
      }
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    if (touched[name] && name === 'name') {
      setErrors(prev => ({ ...prev, name: validateName(value) }));
    }
  };

  const handleBlur = (e) => {
    const { name, value } = e.target;
    setTouched(prev => ({ ...prev, [name]: true }));
    if (name === 'name') {
      setErrors(prev => ({ ...prev, name: validateName(value) }));
    } else if (name === 'price') {
      setErrors(prev => ({ ...prev, price: validatePrice(value) }));
    }
  };

  const isFormValid = () => {
    const nameErrors = validateName(formData.name);
    const priceErrors = validatePrice(formData.price);
    setErrors({ name: nameErrors, price: priceErrors });
    return nameErrors.length === 0 && priceErrors.length === 0;
  };

  // Drag & Drop (без изменений)
  const handleDragOver = (e) => { e.preventDefault(); setIsDragging(true); };
  const handleDragLeave = (e) => { e.preventDefault(); setIsDragging(false); };
  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);
    const droppedFiles = Array.from(e.dataTransfer.files);
    setFiles(prev => [...prev, ...droppedFiles]);
  };
  const handleFileSelect = (e) => {
    const selectedFiles = Array.from(e.target.files);
    setFiles(prev => [...prev, ...selectedFiles]);
  };

  // ← НОВОЕ: Отправка формы с запросом к API
  const handleSubmit = async (e) => {
    e.preventDefault();
    setSubmitError(null);

    // 1. Валидация
    if (!isFormValid()) {
      setTouched({ name: true, price: true });
      return;
    }

    try {
      setIsSubmitting(true);

      // 2. Получаем вишлист (личный если нет eventId, иначе — игровой)
      const wishlistData = await fetchMyWishlist(eventId);
      const wishlistId = wishlistData.id || wishlistData.wishlistId;

      if (!wishlistId) {
        throw new Error('Не удалось получить ID вишлиста. Возможно, вы ещё не вступили в игру.');
      }

      // 3. Загрузка файла (если выбран)
      let imageURL = '';
      if (files[0]) {
        try {
          const uploadResult = await uploadFile(files[0]);
          imageURL = uploadResult?.url || '';
        } catch (uploadErr) {
          console.warn('Не удалось загрузить изображение:', uploadErr);
        }
      }

      // 4. Подготовка данных товара
      const itemData = {
        title: formData.name.trim(),
        price: parseFloat(formData.price),
        ...(formData.link && { link: formData.link }),
        ...(imageURL && { imageURL }),
      };

      // 5. Отправляем товар на сервер
      await addWishlistItem(wishlistId, itemData);

      // 6. Успех -> переходим обратно
      if (eventId) {
        navigate(`/game/${eventId}/wishlist`);
      } else {
        navigate('/wishlist');
      }

    } catch (error) {
      console.error('Ошибка при добавлении товара:', error);
      setSubmitError(error.message || 'Не удалось добавить товар. Попробуйте позже.');
      alert(error.message || 'Ошибка при добавлении товара');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="overlay_wishlist">
      <div className="card_wishlist">
        <h1>Добавление товара</h1>
        
        <form onSubmit={handleSubmit} noValidate>
          <div className="wishlist-content">
            <div className="form-left">
              <div className="form-group">
                <label>Название товара <span className="required">*</span></label>
                <input 
                  type="text" 
                  name="name"
                  placeholder="Введите название" 
                  value={formData.name}
                  onChange={handleChange}
                  onBlur={handleBlur}
                  required 
                  disabled={isSubmitting}
                  className={errors.name.length > 0 && touched.name ? 'input-error' : ''}
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
                <label>Цена <span className="required">*</span></label>
                <input 
                  type="text" 
                  name="price"
                  placeholder="0" 
                  value={formData.price}
                  onChange={handlePriceChange}
                  onBlur={handleBlur}
                  inputMode="decimal"
                  disabled={isSubmitting}
                  className={errors.price.length > 0 && touched.price ? 'input-error' : ''}
                />
                {errors.price.length > 0 && touched.price && (
                  <ul className="error-list">
                    {errors.price.map((err, i) => (
                      <li key={i} className="error-item">• {err}</li>
                    ))}
                  </ul>
                )}
              </div>
              
              <div className="form-group">
                <label>Ссылка на товар</label>
                <input 
                  type="url" 
                  name="link"
                  placeholder="https://..." 
                  value={formData.link}
                  onChange={handleChange}
                  disabled={isSubmitting}
                />
              </div>
            </div>
            
            <div className="form-right">
              <div 
                className={`upload-area ${isDragging ? 'dragover' : ''}`}
                onClick={() => fileInputRef.current?.click()}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                style={{ opacity: isSubmitting ? 0.6 : 1, pointerEvents: isSubmitting ? 'none' : 'auto' }}
              >
                <i className="ti ti-upload" style={{ fontSize: '48px', color: '#44E858' }}></i>
                <div className="upload-text">
                  Чтобы загрузить файл кликните или<br/>перетащите его в эту область
                </div>
                <div className="upload-hint">
                  Можно загрузить не более 1 файла
                </div>
                <input 
                  ref={fileInputRef}
                  type="file" 
                  style={{ display: 'none' }} 
                  accept="image/*" 
                  multiple
                  onChange={handleFileSelect}
                  disabled={isSubmitting}
                />
              </div>
              
              {files.length > 0 && (
                <div className="file-list">
                  <strong>Выбрано файлов: {files.length}</strong>
                  {files.map((file, index) => (
                    <div key={index} className="file-item">• {file.name}</div>
                  ))}
                </div>
              )}
              
              {submitError && (
                <p style={{ color: '#e74c3c', fontSize: '14px', marginTop: '10px' }}>
                   {submitError}
                </p>
              )}
            </div>
          </div>
          
          <div className="wishlist-add-buttons">
            <button 
              type="submit" 
              className="btn-primary" 
              disabled={isSubmitting}
              style={{ opacity: isSubmitting ? 0.7 : 1, cursor: isSubmitting ? 'not-allowed' : 'pointer' }}
            >
              {isSubmitting ? 'Добавление...' : 'Добавь в вишлист'}
            </button>
            <button 
              type="button" 
              className="btn-secondary"
              onClick={handleGoBack}
              disabled={isSubmitting}
            >
              Отмена
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default WishlistAdd;