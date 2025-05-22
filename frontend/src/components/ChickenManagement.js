import React, { useState, useEffect } from 'react';
import axios from 'axios';

// Настройка базового URL для axios
const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

const ChickenManagement = () => {
  const [chickens, setChickens] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingChicken, setEditingChicken] = useState(null);
  const [formData, setFormData] = useState({
    cage_id: '',
    weight: '',
    age: '',
    egg_per_month: '',
    breed: ''
  });

  useEffect(() => {
    fetchChickens();
  }, []);

  const fetchChickens = async () => {
    try {
      setLoading(true);
      const response = await api.get('/api/chickens');
      setChickens(response.data || []);
    } catch (err) {
      setError('Ошибка при загрузке списка кур');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    
    try {
      const data = {
        cage_id: parseInt(formData.cage_id),
        weight: parseFloat(formData.weight),
        age: parseInt(formData.age),
        egg_per_month: parseInt(formData.egg_per_month),
        breed: formData.breed.trim()
      };

      console.log('Отправляемые данные:', data);

      if (editingChicken) {
        await api.put(`/api/chickens/${editingChicken.id}`, data);
        setSuccess('Курица успешно обновлена');
      } else {
        await api.post('/api/chickens', data);
        setSuccess('Курица успешно добавлена');
      }

      resetForm();
      fetchChickens();
    } catch (err) {
      console.error('Ошибка запроса:', err);
      setError(err.response?.data?.error || 'Ошибка при сохранении курицы');
    }
  };

  const handleEdit = (chicken) => {
    setEditingChicken(chicken);
    setFormData({
      cage_id: chicken.cage_id.toString(),
      weight: chicken.weight.toString(),
      age: chicken.age.toString(),
      egg_per_month: chicken.egg_per_month.toString(),
      breed: chicken.breed
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить эту курицу?')) {
      try {
        await api.delete(`/api/chickens/${id}`);
        setSuccess('Курица успешно удалена');
        fetchChickens();
      } catch (err) {
        console.error('Ошибка удаления:', err);
        setError('Ошибка при удалении курицы');
      }
    }
  };

  const resetForm = () => {
    setFormData({
      cage_id: '',
      weight: '',
      age: '',
      egg_per_month: '',
      breed: ''
    });
    setEditingChicken(null);
    setShowModal(false);
  };

  const handleInputChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  return (
    <div className="page">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1 className="page-title">Управление курами</h1>
        <button 
          className="btn btn-primary"
          onClick={() => setShowModal(true)}
        >
          Добавить курицу
        </button>
      </div>

      {error && <div className="error">{error}</div>}
      {success && <div className="success">{success}</div>}

      {loading ? (
        <div className="loading">Загрузка...</div>
      ) : (
        <div className="table-container">
          <table className="table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Клетка</th>
                <th>Порода</th>
                <th>Вес (кг)</th>
                <th>Возраст (мес.)</th>
                <th>Яиц/месяц</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {chickens.length === 0 ? (
                <tr>
                  <td colSpan="7" style={{ textAlign: 'center', padding: '2rem' }}>
                    Пока нет ни одной курицы
                  </td>
                </tr>
              ) : (
                chickens.map(chicken => (
                  <tr key={chicken.id}>
                    <td>{chicken.id}</td>
                    <td>#{chicken.cage_id}</td>
                    <td>{chicken.breed}</td>
                    <td>{chicken.weight}</td>
                    <td>{chicken.age}</td>
                    <td>{chicken.egg_per_month}</td>
                    <td>
                      <div className="table-actions">
                        <button 
                          className="btn btn-warning btn-small"
                          onClick={() => handleEdit(chicken)}
                        >
                          Изменить
                        </button>
                        <button 
                          className="btn btn-danger btn-small"
                          onClick={() => handleDelete(chicken.id)}
                        >
                          Удалить
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      )}

      {/* Модальное окно для добавления/редактирования */}
      {showModal && (
        <div className="modal">
          <div className="modal-content">
            <div className="modal-header">
              <h2 className="modal-title">
                {editingChicken ? 'Редактировать курицу' : 'Добавить курицу'}
              </h2>
            </div>
            
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label className="form-label">Номер клетки:</label>
                <input
                  type="number"
                  name="cage_id"
                  className="form-input"
                  value={formData.cage_id}
                  onChange={handleInputChange}
                  required
                  min="1"
                />
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Вес (кг):</label>
                  <input
                    type="number"
                    name="weight"
                    className="form-input"
                    value={formData.weight}
                    onChange={handleInputChange}
                    required
                    min="0"
                    step="0.1"
                    placeholder="Например: 2.5"
                  />
                </div>

                <div className="form-group">
                  <label className="form-label">Возраст (месяцы):</label>
                  <input
                    type="number"
                    name="age"
                    className="form-input"
                    value={formData.age}
                    onChange={handleInputChange}
                    required
                    min="1"
                    placeholder="Например: 12"
                  />
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Яиц в месяц:</label>
                  <input
                    type="number"
                    name="egg_per_month"
                    className="form-input"
                    value={formData.egg_per_month}
                    onChange={handleInputChange}
                    required
                    min="0"
                    placeholder="Например: 25"
                  />
                </div>

                <div className="form-group">
                  <label className="form-label">Порода:</label>
                  <input
                    type="text"
                    name="breed"
                    className="form-input"
                    value={formData.breed}
                    onChange={handleInputChange}
                    required
                    placeholder="Например: Леггорн"
                  />
                </div>
              </div>

              <div className="form-actions">
                <button type="submit" className="btn btn-primary">
                  {editingChicken ? 'Обновить' : 'Добавить'}
                </button>
                <button 
                  type="button" 
                  className="btn btn-secondary"
                  onClick={resetForm}
                >
                  Отмена
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default ChickenManagement;