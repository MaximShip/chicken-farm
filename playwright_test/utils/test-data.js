export const testData = {
  chickens: {
    valid: {
      cage_id: 6,
      weight: 2.8,
      age: 15,
      egg_per_month: 28,
      breed: 'Нью-Гемпшир E2E'
    },
    validAlternative: {
      cage_id: 7,
      weight: 3.2,
      age: 20,
      egg_per_month: 30,
      breed: 'Кучинская Юбилейная'
    },
    update: {
      cage_id: 6,
      weight: 3.0,
      age: 16,
      egg_per_month: 32,
      breed: 'Нью-Гемпшир E2E Updated'
    },
    invalid: {
      empty: {
        cage_id: '',
        weight: '',
        age: '',
        egg_per_month: '',
        breed: ''
      },
      negative: {
        cage_id: 1,
        weight: -1,
        age: -5,
        egg_per_month: -10,
        breed: 'Invalid'
      }
    }
  },

  employees: {
    valid: {
      full_name: 'Сидоров Сидор Сидорович',
      passport_data: '3456 789012',
      salary: 55000,
      cages: '4, 5'
    },
    validAlternative: {
      full_name: 'Козлов Козлов Козлович',
      passport_data: '4567 890123',
      salary: 48000,
      cages: '6'
    },
    update: {
      full_name: 'Сидоров Сидор Сидорович Updated',
      passport_data: '3456 789012',
      salary: 62000,
      cages: '4, 5, 6'
    },
    invalid: {
      empty: {
        full_name: '',
        passport_data: '',
        salary: '',
        cages: ''
      },
      wrongPassport: {
        full_name: 'Неверный Паспорт',
        passport_data: '123456789', 
        salary: 50000,
        cages: '1'
      },
      negativeSalary: {
        full_name: 'Отрицательная Зарплата',
        passport_data: '1111 222222',
        salary: -1000,
        cages: '1'
      }
    }
  },

  dateRanges: {
    currentMonth: {
      start: '2025-01-01',
      end: '2025-01-31'
    },
    lastMonth: {
      start: '2024-12-01',
      end: '2024-12-31'
    },
    quarter: {
      start: '2025-01-01',
      end: '2025-03-31'
    }
  }
};

export const generators = {
  randomChicken: () => ({
    cage_id: Math.floor(Math.random() * 100) + 10,
    weight: (Math.random() * 3 + 1.5).toFixed(1),
    age: Math.floor(Math.random() * 24) + 6,
    egg_per_month: Math.floor(Math.random() * 20) + 15,
    breed: `Random Breed ${Math.floor(Math.random() * 1000)}`
  }),

  randomEmployee: () => {
    const firstNames = ['Иван', 'Петр', 'Сидор', 'Алексей', 'Николай'];
    const lastNames = ['Иванов', 'Петров', 'Сидоров', 'Алексеев', 'Николаев'];
    const middleNames = ['Иванович', 'Петрович', 'Сидорович', 'Алексеевич', 'Николаевич'];
    
    const firstName = firstNames[Math.floor(Math.random() * firstNames.length)];
    const lastName = lastNames[Math.floor(Math.random() * lastNames.length)];
    const middleName = middleNames[Math.floor(Math.random() * middleNames.length)];
    
    const passportSeries = Math.floor(Math.random() * 9000) + 1000;
    const passportNumber = Math.floor(Math.random() * 900000) + 100000;
    
    return {
      full_name: `${lastName} ${firstName} ${middleName}`,
      passport_data: `${passportSeries} ${passportNumber}`,
      salary: Math.floor(Math.random() * 50000) + 30000,
      cages: `${Math.floor(Math.random() * 5) + 1}, ${Math.floor(Math.random() * 5) + 6}`
    };
  }
};
