import React, { useState, useEffect } from 'react';
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

const EmployeeManagement = () => {
  const [employees, setEmployees] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState(null);
  const [formData, setFormData] = useState({
    full_name: '',
    passport_data: '',
    salary: '',
    cages: ''
  });

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      const response = await api.get('/api/employees');
      setEmployees(response.data || []);
    } catch (err) {
      setError('Ошибка при загрузке списка работников');
      console.error('Ошибка загрузки работников:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    
    try {
      const cagesArray = formData.cages
        .split(',')
        .map(id => parseInt(id.trim()))
        .filter(id => !isNaN(id) && id > 0);

      const data = {
        full_name: formData.full_name.trim(),
        passport_data: formData.passport_data.trim(),
        salary: parseFloat(formData.salary),
        cages: cagesArray
      };

      console.log('Отправляемые данные работника:', data);

      if (editingEmployee) {
        await api.put(`/api/employees/${editingEmployee.id}`, data);
        setSuccess('Работник успешно обновлен');
      } else {
        await api.post('/api/employees', data);
        setSuccess('Работник успешно добавлен');
      }

      resetForm();
      fetchEmployees();
    } catch (err) {
      console.error('Ошибка сохранения работника:', err);
      setError(err.response?.data?.error || 'Ошибка при сохранении работника');
    }
  };

  const handleEdit = (employee) => {
    setEditingEmployee(employee);
    setFormData({
      full_name: employee.full_name,
      passport_data: employee.passport_data,
      salary: employee.salary.toString(),
      cages: employee.cages ? employee.cages.join(', ') : ''
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить этого работника?')) {
      try {
        await api.delete(`/api/employees/${id}`);
        setSuccess('Работник успешно удален');
        fetchEmployees();
      } catch (err) {
        console.error('Ошибка удаления работника:', err);
        setError('Ошибка при удалении работника');
      }
    }
  };

  const resetForm = () => {
    setFormData({
      full_name: '',
      passport_data: '',
      salary: '',
      cages: ''
    });
    setEditingEmployee(null);
    setShowModal(false);
  };

  const handleInputChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const validatePassportData = (passport) => {
    const passportRegex = /^\d{4}\s\d{6}$/;
    return passportRegex.test(passport);
  };

  const handleFormSubmit = (e) => {
    e.preventDefault();
    
    if (!validatePassportData(formData.passport_data.trim())) {
      setError('Паспортные данные должны быть в формате: 1234 567890');
      return;
    }

    if (parseFloat(formData.salary) <= 0) {
      setError('Зарплата должна быть больше нуля');
      return;
    }

    handleSubmit(e);
  };

  return (
    <div className="page">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
        <h1 className="page-title">Управление работниками</h1>
        <button 
          className="btn btn-primary"
          onClick={() => setShowModal(true)}
        >
          Добавить работника
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
                <th>ФИО</th>
                <th>Паспорт</th>
                <th>Зарплата</th>
                <th>Клетки</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {employees.length === 0 ? (
                <tr>
                  <td colSpan="6" style={{ textAlign: 'center', padding: '2rem' }}>
                    Пока нет ни одного работника
                  </td>
                </tr>
              ) : (
                employees.map(employee => (
                  <tr key={employee.id}>
                    <td>{employee.id}</td>
                    <td>{employee.full_name}</td>
                    <td>{employee.passport_data}</td>
                    <td>{employee.salary.toLocaleString()} ₽</td>
                    <td>
                      {employee.cages && employee.cages.length > 0
                        ? employee.cages.map(id => `#${id}`).join(', ')
                        : 'Не назначены'
                      }
                    </td>
                    <td>
                      <div className="table-actions">
                        <button 
                          className="btn btn-warning btn-small"
                          onClick={() => handleEdit(employee)}
                        >
                          Изменить
                        </button>
                        <button 
                          className="btn btn-danger btn-small"
                          onClick={() => handleDelete(employee.id)}
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

      {showModal && (
        <div className="modal">
          <div className="modal-content">
            <div className="modal-header">
              <h2 className="modal-title">
                {editingEmployee ? 'Редактировать работника' : 'Добавить работника'}
              </h2>
            </div>
            
            <form onSubmit={handleFormSubmit}>
              <div className="form-group">
                <label className="form-label">ФИО:</label>
                <input
                  type="text"
                  name="full_name"
                  className="form-input"
                  value={formData.full_name}
                  onChange={handleInputChange}
                  required
                  placeholder="Иванов Иван Иванович"
                />
              </div>

              <div className="form-group">
                <label className="form-label">Паспортные данные:</label>
                <input
                  type="text"
                  name="passport_data"
                  className="form-input"
                  value={formData.passport_data}
                  onChange={handleInputChange}
                  required
                  placeholder="1234 567890"
                  pattern="\d{4}\s\d{6}"
                  title="Формат: 1234 567890"
                />
                <small style={{ color: '#7f8c8d', fontSize: '0.875rem' }}>
                  Формат: 1234 567890 (серия пробел номер)
                </small>
              </div>

              <div className="form-group">
                <label className="form-label">Зарплата (руб.):</label>
                <input
                  type="number"
                  step="0.01"
                  name="salary"
                  className="form-input"
                  value={formData.salary}
                  onChange={handleInputChange}
                  required
                  min="0.01"
                  placeholder="50000"
                />
              </div>

              <div className="form-group">
                <label className="form-label">Клетки (через запятую):</label>
                <input
                  type="text"
                  name="cages"
                  className="form-input"
                  value={formData.cages}
                  onChange={handleInputChange}
                  placeholder="1, 2, 3"
                />
                <small style={{ color: '#7f8c8d', fontSize: '0.875rem' }}>
                  Укажите номера клеток через запятую. Поле необязательное.
                </small>
              </div>

              <div className="form-actions">
                <button type="submit" className="btn btn-primary">
                  {editingEmployee ? 'Обновить' : 'Добавить'}
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

export default EmployeeManagement;