import React from 'react';
import styled from 'styled-components';
import { FiExternalLink, FiBookmark, FiShare2 } from 'react-icons/fi';

const DealsContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const DealItem = styled.div`
  display: flex;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:last-child {
    border-bottom: none;
  }
`;

const DealImage = styled.img`
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
`;

const DealContent = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const DealTitle = styled.h4`
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: white;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const DealDescription = styled.p`
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const DealPrice = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 5px 0;
`;

const CurrentPrice = styled.span`
  font-size: 16px;
  font-weight: 700;
  color: #4caf50;
`;

const OriginalPrice = styled.span`
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  text-decoration: line-through;
`;

const Discount = styled.span`
  font-size: 10px;
  background: #ff6b6b;
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
`;

const DealMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
`;

const DealActions = styled.div`
  display: flex;
  gap: 8px;
`;

const ActionButton = styled.button`
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s ease;

  &:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
  }
`;

const LoadingMessage = styled.div`
  text-align: center;
  padding: 40px;
  color: rgba(255, 255, 255, 0.6);
`;

const StaticDataWarning = styled.div`
  background: rgba(255, 193, 7, 0.1);
  border: 1px solid rgba(255, 193, 7, 0.3);
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  color: #ffc107;
  font-size: 14px;
  
  strong {
    color: #ff9800;
  }
  
  a {
    color: #4fc3f7;
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
`;

const ApiStatusIndicator = styled.div`
  display: inline-block;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  margin-left: 8px;
  
  &.real-api {
    background: rgba(76, 175, 80, 0.2);
    color: #4caf50;
    border: 1px solid rgba(76, 175, 80, 0.3);
  }
  
  &.static-data {
    background: rgba(255, 107, 107, 0.2);
    color: #ff6b6b;
    border: 1px solid rgba(255, 107, 107, 0.3);
  }
`;

function DealsSection({ data, loading, onUserAction }) {
  if (loading) {
    return <LoadingMessage>Loading amazing deals...</LoadingMessage>;
  }

  if (!data || !data.deals) {
    return <LoadingMessage>No deals available at the moment.</LoadingMessage>;
  }

  // Check if this is static data
  const isStaticData = data.error === "STATIC DATA - No Amazon API Key Found" || 
                      data.deals.some(deal => deal.is_static) ||
                      data.source !== "REAL AMAZON API DATA";

  const handleAction = (action, deal) => {
    onUserAction(action, deal.id, 'deals');
  };

  return (
    <DealsContainer>
      {isStaticData && (
        <StaticDataWarning>
          ⚠️ <strong>STATIC DATA WARNING:</strong> {data.message || 'This is mock data. Add AMAZON_API_KEY to env.local for real deals.'}
          <br />
          <small>Get API key: <a href="https://webservices.amazon.com/" target="_blank" rel="noopener noreferrer">https://webservices.amazon.com/</a></small>
        </StaticDataWarning>
      )}
      
      {!isStaticData && (
        <div style={{ marginBottom: '15px', textAlign: 'center' }}>
          <ApiStatusIndicator className="real-api">Real Amazon API</ApiStatusIndicator>
        </div>
      )}
      
      {data.deals.slice(0, 4).map((deal, index) => (
        <DealItem key={deal.id || index}>
          <DealImage 
            src={deal.image_url || '/placeholder-deal.jpg'} 
            alt={deal.title}
            onError={(e) => {
              e.target.style.display = 'none';
            }}
          />
          <DealContent>
            <DealTitle>{deal.title}</DealTitle>
            <DealDescription>
              {deal.description?.substring(0, 100)}...
            </DealDescription>
            <DealPrice>
              <CurrentPrice>${deal.current_price}</CurrentPrice>
              {deal.original_price && (
                <OriginalPrice>${deal.original_price}</OriginalPrice>
              )}
              {deal.discount && (
                <Discount>{deal.discount}% OFF</Discount>
              )}
            </DealPrice>
            <DealMeta>
              <span>{deal.store} • {deal.category}</span>
              <DealActions>
                <ActionButton 
                  onClick={() => handleAction('bookmark', deal)}
                  title="Save Deal"
                >
                  <FiBookmark />
                </ActionButton>
                <ActionButton 
                  onClick={() => handleAction('share', deal)}
                  title="Share"
                >
                  <FiShare2 />
                </ActionButton>
                <ActionButton 
                  onClick={() => {
                    if (deal.is_static) {
                      // For static data, redirect to API registration
                      window.open('https://webservices.amazon.com/', '_blank', 'noopener,noreferrer');
                    } else if (deal.url && deal.url !== 'https://example.com/deals/1') {
                      // Real deal
                      window.open(deal.url, '_blank', 'noopener,noreferrer');
                    } else {
                      // Fallback
                      window.open('https://amazon.com', '_blank', 'noopener,noreferrer');
                    }
                    handleAction('click', deal);
                  }}
                  title={deal.is_static ? "Get API Key" : "View Deal"}
                >
                  <FiExternalLink />
                </ActionButton>
              </DealActions>
            </DealMeta>
          </DealContent>
        </DealItem>
      ))}
    </DealsContainer>
  );
}

export default DealsSection;