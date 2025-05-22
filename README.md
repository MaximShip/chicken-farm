## Backend

    - cd backend
    - go mod tidy
    - go run main.go

## Frontend

    - npm install
    - npm start


## Testing 

    
    - Backend 
        - cd backend/cmd
        - go test main_test.go

    - Frontend (Playwright)
        - npm install --save-dev @playwright/test
        - npx playwright install chromium
        - npm run test
