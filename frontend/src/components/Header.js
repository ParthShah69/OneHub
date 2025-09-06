import React from 'react';
import styled from 'styled-components';
import { FiUser, FiLogOut } from 'react-icons/fi';

const HeaderContainer = styled.header`
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 15px 0;
`;

const HeaderContent = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Logo = styled.div`
  font-size: 24px;
  font-weight: 700;
  color: white;
  display: flex;
  align-items: center;
  gap: 10px;
`;

const UserSection = styled.div`
  display: flex;
  align-items: center;
  gap: 15px;
`;

const UserInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  color: white;
  font-size: 14px;
`;

const LogoutButton = styled.button`
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  color: white;
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  transition: all 0.3s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.2);
    border-color: rgba(255, 255, 255, 0.5);
  }
`;

function Header({ user }) {
  const handleLogout = () => {
    localStorage.removeItem('dashboard_user');
    window.location.reload();
  };

  return (
    <HeaderContainer>
      <HeaderContent>
        <Logo>
          🚀 Personalized Dashboard
        </Logo>
        {user && (
          <UserSection>
            <UserInfo>
              <FiUser />
              {user.name}
            </UserInfo>
            <LogoutButton onClick={handleLogout}>
              <FiLogOut />
              Logout
            </LogoutButton>
          </UserSection>
        )}
      </HeaderContent>
    </HeaderContainer>
  );
}

export default Header;
