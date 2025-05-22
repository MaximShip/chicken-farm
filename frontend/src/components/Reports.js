import React, { useState, useEffect } from 'react';
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

const Reports = () => {
  const [eggStats, setEggStats] = useState(null);
  const [employeeStats, setEmployeeStats] = useState([]);
  const [lowProductivityChickens, setLowProductivityChickens] = useState([]);
  const [mostProductiveChicken, setMostProductiveChicken] = useState(null);
  const [employeeChickenCounts, setEmployeeChickenCounts] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [dateRange, setDateRange] = useState({
    start_date: '2025-01-01',
    end_date: '2025-01-31'
  });

  useEffect(() => {
    fetchStaticReports();
  }, []);

  const fetchStaticReports = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const [lowProdRes, mostProdRes, empCountRes] = await Promise.all([
        api.get('/api/reports/low-productivity-chickens').catch((err) => {
          console.log('Ошибка загрузки малопродуктивных кур:', err);
          return { data: [] };
        }),
        api.get('/api/reports/most-productive-chicken').catch((err) => {
          console.log('Ошибка загрузки самой продуктивной курицы:', err);
          return { data: null };
        }),
        api.get('/api/reports/employee-chicken-counts').catch((err) => {
          console.log('Ошибка загрузки статистики по работникам:', err);
          return { data: [] };
        })
      ]);

      setLowProductivityChickens(lowProdRes.data || []);
      setMostProductiveChicken(mostProdRes.data);
      setEmployeeChickenCounts(empCountRes.data || []);

      console.log('Статические отчеты загружены:', {
        lowProductivity: lowProdRes.data?.length || 0,
        mostProductive: !!mostProdRes.data,
        employeeCounts: empCountRes.data?.length || 0
      });

    } catch (err) {
      console.error('Ошибка загрузки статических отчетов:', err);
      setError('Ошибка при загрузке отчетов');
    } finally {
      setLoading(false);
    }
  };

  const fetchDateRangeReports = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log('Загрузка отчетов за период:', dateRange);

      const [eggStatsRes, empStatsRes] = await Promise.all([
        api.get('/api/reports/egg-stats', { params: dateRange }),
        api.get('/api/reports/employee-egg-stats', { params: dateRange })
      ]);

      setEggStats(eggStatsRes.data);
      setEmployeeStats(empStatsRes.data.stats || []);

      console.log('Отчеты за период загружены:', {
        eggStats: eggStatsRes.data,
        employeeStats: empStatsRes.data.stats?.length || 0
      });

    } catch (err) {
      console.error('Ошибка загрузки отчетов за период:', err);
      setError('Ошибка при загрузке отчетов за период: ' + (err.response?.data?.error || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleDateChange = (e) => {
    setDateRange({
      ...dateRange,
      [e.target.name]: e.target.value
    });
  };

  const getCurrentDate = () => {
    return new Date().toISOString().split('T')[0];
  };

  const getDateMonthAgo = () => {
    const date = new Date();
    date.setMonth(date.getMonth() - 1);
    return date.toISOString().split('T')[0];
  };

  const setPresetPeriod = (period) => {
    const today = getCurrentDate();
    let startDate = today;

    switch (period) {
      case 'week':
        const weekAgo = new Date();
        weekAgo.setDate(weekAgo.getDate() - 7);
        startDate = weekAgo.toISOString().split('T')[0];
        break;
      case 'month':
        startDate = getDateMonthAgo();
        break;
      case 'year':
        const yearAgo = new Date();
        yearAgo.setFullYear(yearAgo.getFullYear() - 1);
        startDate = yearAgo.toISOString().split('T')[0];
        break;
      default:
        startDate = today;
    }

    setDateRange({
      start_date: startDate,
      end_date: today
    });
  };

  return (
    <div className="page">
      <h1 className="page-title">Отчеты и аналитика</h1>

      {error && <div className="error">{error}</div>}

      <div className="page-section" style={{ marginBottom: '2rem' }}>
        <h2 className="page-subtitle">Отчеты за период</h2>
        
        <div style={{ marginBottom: '1rem' }}>
          <span style={{ marginRight: '1rem', fontWeight: 'bold' }}>Быстрый выбор:</span>
          <button 
            className="btn btn-secondary btn-small" 
            onClick={() => setPresetPeriod('week')}
            style={{ marginRight: '0.5rem' }}
          >
            Неделя
          </button>
          <button 
            className="btn btn-secondary btn-small" 
            onClick={() => setPresetPeriod('month')}
            style={{ marginRight: '0.5rem' }}
          >
            Месяц
          </button>
          <button 
            className="btn btn-secondary btn-small" 
            onClick={() => setPresetPeriod('year')}
          >
            Год
          </button>
        </div>

        <div style={{ display: 'flex', gap: '1rem', alignItems: 'end', flexWrap: 'wrap' }}>
          <div className="form-group" style={{ marginBottom: 0 }}>
            <label className="form-label">Начальная дата:</label>
            <input
              type="date"
              name="start_date"
              className="form-input"
              value={dateRange.start_date}
              onChange={handleDateChange}
              style={{ width: 'auto' }}
              max={getCurrentDate()}
            />
          </div>
          <div className="form-group" style={{ marginBottom: 0 }}>
            <label className="form-label">Конечная дата:</label>
            <input
              type="date"
              name="end_date"
              className="form-input"
              value={dateRange.end_date}
              onChange={handleDateChange}
              style={{ width: 'auto' }}
              min={dateRange.start_date}
              max={getCurrentDate()}
            />
          </div>
          <button 
            className="btn btn-primary"
            onClick={fetchDateRangeReports}
            disabled={loading || !dateRange.start_date || !dateRange.end_date}
          >
            {loading ? 'Загрузка...' : 'Загрузить отчеты'}
          </button>
        </div>
      </div>

      {eggStats && (
        <div className="page-section">
          <h2 className="page-subtitle">Статистика по яйцам за период</h2>
          <div className="stats-grid">
            <div className="stat-card">
              <div className="stat-value">{eggStats.total_eggs || 0}</div>
              <div className="stat-label">Общее количество яиц</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">{(eggStats.total_cost || 0).toFixed(2)} ₽</div>
              <div className="stat-label">Общая стоимость</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">
                {eggStats.total_eggs > 0 ? (eggStats.total_cost / eggStats.total_eggs).toFixed(2) : '0.00'} ₽
              </div>
              <div className="stat-label">Цена за яйцо</div>
            </div>
          </div>
        </div>
      )}

      {employeeStats.length > 0 && (
        <div className="page-section">
          <h2 className="page-subtitle">Количество яиц, собранных работниками</h2>
          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Работник</th>
                  <th>Количество яиц</th>
                  <th>% от общего количества</th>
                </tr>
              </thead>
              <tbody>
                {employeeStats.map(stat => {
                  const totalEggs = employeeStats.reduce((sum, s) => sum + s.egg_count, 0);
                  const percentage = totalEggs > 0 ? ((stat.egg_count / totalEggs) * 100).toFixed(1) : '0';
                  
                  return (
                    <tr key={stat.employee_id}>
                      <td>{stat.employee_name}</td>
                      <td>{stat.egg_count}</td>
                      <td>{percentage}%</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}

      <div className="page-section">
        <h2 className="page-subtitle">Куры с низкой продуктивностью</h2>
        {lowProductivityChickens.length > 0 ? (
          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Клетка</th>
                  <th>Порода</th>
                  <th>Яиц/месяц</th>
                  <th>Вес (кг)</th>
                  <th>Возраст (мес.)</th>
                </tr>
              </thead>
              <tbody>
                {lowProductivityChickens.map(chicken => (
                  <tr key={chicken.id}>
                    <td>{chicken.id}</td>
                    <td>#{chicken.cage_id}</td>
                    <td>{chicken.breed}</td>
                    <td className="text-warning">{chicken.egg_per_month}</td>
                    <td>{chicken.weight}</td>
                    <td>{chicken.age}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div style={{ textAlign: 'center', padding: '2rem', color: '#27ae60' }}>
            ✓ Все куры имеют продуктивность выше средней
          </div>
        )}
      </div>

      {mostProductiveChicken && (
        <div className="page-section">
          <h2 className="page-subtitle">Самая продуктивная кура</h2>
          <div className="stat-card" style={{ maxWidth: '400px' }}>
            <div className="stat-value">Клетка #{mostProductiveChicken.cage_number}</div>
            <div className="stat-label">
              <div>ID курицы: {mostProductiveChicken.chicken_id}</div>
              <div>Яиц в месяц: {mostProductiveChicken.egg_per_month}</div>
            </div>
          </div>
        </div>
      )}

      <div className="page-section">
        <h2 className="page-subtitle">Количество кур, обслуживаемых работниками</h2>
        {employeeChickenCounts.length > 0 ? (
          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Работник</th>
                  <th>Количество кур</th>
                  <th>Нагрузка</th>
                </tr>
              </thead>
              <tbody>
                {employeeChickenCounts.map(stat => {
                  const workload = stat.chicken_count === 0 ? 'Нет' : 
                    stat.chicken_count <= 2 ? 'Низкая' :
                    stat.chicken_count <= 4 ? 'Средняя' : 'Высокая';
                  
                  return (
                    <tr key={stat.employee_id}>
                      <td>{stat.employee_name}</td>
                      <td>{stat.chicken_count}</td>
                      <td>{workload}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        ) : (
          <div style={{ textAlign: 'center', padding: '2rem', color: '#7f8c8d' }}>
            Нет данных о работниках
          </div>
        )}
      </div>

      <div className="page-section">
        <h2 className="page-subtitle">Как использовать отчеты</h2>
        <ul>
          <li><strong>Отчеты за период:</strong> Выберите даты и нажмите "Загрузить отчеты" для получения статистики по яйцам и работникам за указанный период</li>
          <li><strong>Быстрый выбор периода:</strong> Используйте кнопки для выбора стандартных периодов (неделя, месяц, год)</li>
          <li><strong>Куры с низкой продуктивностью:</strong> Автоматически показывает кур, чья продуктивность ниже средней по фабрике</li>
          <li><strong>Самая продуктивная кура:</strong> Отображает курицу с наибольшим количеством яиц в месяц</li>
          <li><strong>Количество кур по работникам:</strong> Показывает нагрузку каждого работника и помогает распределить обязанности</li>
        </ul>
      </div>
    </div>
  );
};

export default Reports;