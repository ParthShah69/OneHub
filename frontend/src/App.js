import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import axios from 'axios';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Dashboard from './components/Dashboard';
import Onboarding from './components/Onboarding';
import Header from './components/Header';
import './App.css';

const AppContainer = styled.div`
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
`;

const MainContent = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
`;

function App() {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if user exists in localStorage
    const savedUser = localStorage.getItem('dashboard_user');
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
    setIsLoading(false);
  }, []);

  const handleUserCreate = (userData) => {
    setUser(userData);
    localStorage.setItem('dashboard_user', JSON.stringify(userData));
  };

  if (isLoading) {
    return (
      <AppContainer>
        <MainContent>
          <div style={{ textAlign: 'center', padding: '50px', color: 'white' }}>
            <h2>Loading your personalized dashboard...</h2>
          </div>
        </MainContent>
      </AppContainer>
    );
  }

  return (
    <Router>
      <AppContainer>
        <Header user={user} />
        <MainContent>
          <Routes>
            <Route 
              path="/" 
              element={
                user ? 
                <Dashboard user={user} /> : 
                <Onboarding onUserCreate={handleUserCreate} />
              } 
            />
            <Route 
              path="/dashboard" 
              element={<Dashboard user={user} />} 
            />
          </Routes>
        </MainContent>
      </AppContainer>
    </Router>
  );
}

export default App;
