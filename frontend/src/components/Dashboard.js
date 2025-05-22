import React, { useState, useEffect } from 'react';
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

const Dashboard = () => {
  const [stats, setStats] = useState({
    totalChickens: 0,
    totalEmployees: 0,
    avgEggProduction: 0,
    mostProductiveChicken: null
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      
      const [chickensRes, employeesRes, mostProductiveRes] = await Promise.all([
        api.get('/api/chickens'),
        api.get('/api/employees'),
        api.get('/api/chickens/most-productive').catch((err) => {
          console.log('Ошибка получения самой продуктивной курицы:', err);
          return { data: null };
        })
      ]);

      const chickens = chickensRes.data || [];
      const employees = employeesRes.data || [];
      const avgProduction = chickens.length > 0 
        ? chickens.reduce((sum, chicken) => sum + chicken.egg_per_month, 0) / chickens.length
        : 0;

      setStats({
        totalChickens: chickens.length,
        totalEmployees: employees.length,
        avgEggProduction: Math.round(avgProduction * 10) / 10,
        mostProductiveChicken: mostProductiveRes.data
      });

      console.log('Статистика загружена:', {
        chickens: chickens.length,
        employees: employees.length,
        avgProduction: avgProduction
      });

    } catch (err) {
      console.error('Ошибка загрузки дашборда:', err);
      setError('Ошибка при загрузке данных дашборда');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="page">
        <div className="loading">Загрузка статистики...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="page">
        <div className="error">{error}</div>
        <button 
          className="btn btn-primary" 
          onClick={fetchDashboardData}
          style={{ marginTop: '1rem' }}
        >
          Повторить попытку
        </button>
      </div>
    );
  }

  return (
    <div className="page">
      <h1 className="page-title">Главная панель</h1>
      
      {/* Статистические карточки */}
      <div className="stats-grid">
        <div className="stat-card">
          <div className="stat-value">{stats.totalChickens}</div>
          <div className="stat-label">Всего кур</div>
        </div>
        
        <div className="stat-card">
          <div className="stat-value">{stats.totalEmployees}</div>
          <div className="stat-label">Всего работников</div>
        </div>
        
        <div className="stat-card">
          <div className="stat-value">{stats.avgEggProduction}</div>
          <div className="stat-label">Среднее яиц/месяц</div>
        </div>
        
        {stats.mostProductiveChicken && (
          <div className="stat-card">
            <div className="stat-value">{stats.mostProductiveChicken.egg_per_month}</div>
            <div className="stat-label">Макс. продуктивность</div>
          </div>
        )}
      </div>

      {/* Быстрые действия */}
      <div className="page-section">
        <h2 className="page-subtitle">Быстрые действия</h2>
        <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap' }}>
          <a href="/chickens" className="btn btn-primary">
            Управление курами
          </a>
          <a href="/employees" className="btn btn-secondary">
            Управление работниками
          </a>
          <a href="/reports" className="btn btn-success">
            Просмотр отчетов
          </a>
        </div>
      </div>

      {/* Информация о системе */}
      <div className="page-section">
        <h2 className="page-subtitle">О системе</h2>
        <p>
          Система управления птицефабрикой позволяет вести учет поголовья кур, 
          управлять персоналом и формировать аналитические отчеты. 
          Система разработана как часть курсовой работы по дисциплине 
          "Создание программного обеспечения".
        </p>
        
        <h3>Основные возможности:</h3>
        <ul>
          <li>Ведение базы данных кур с информацией о весе, возрасте и продуктивности</li>
          <li>Управление персоналом и распределение обязанностей</li>
          <li>Формирование отчетов по производительности</li>
          <li>Анализ эффективности работы фабрики</li>
        </ul>
      </div>
    </div>
  );
};

export default Dashboard;