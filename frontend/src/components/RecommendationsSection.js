import React from 'react';
import styled from 'styled-components';
import { FiStar, FiExternalLink, FiBookmark } from 'react-icons/fi';

const RecommendationsContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const RecommendationItem = styled.div`
  display: flex;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:last-child {
    border-bottom: none;
  }
`;

const RecommendationImage = styled.img`
  width: 80px;
  height: 60px;
  object-fit: cover;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
`;

const RecommendationContent = styled.div`
  flex: 1;
`;

const RecommendationTitle = styled.h4`
  margin: 0 0 8px 0;
  font-size: 14px;
  line-height: 1.4;
  color: white;
  display: flex;
  align-items: center;
  gap: 8px;
`;

const RecommendationDescription = styled.p`
  margin: 0 0 10px 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.3;
`;

const RecommendationMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
`;

const RecommendationScore = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  color: #ffd700;
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

const ErrorMessage = styled.div`
  text-align: center;
  padding: 20px;
  color: #ff6b6b;
  background: rgba(255, 107, 107, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(255, 107, 107, 0.2);
`;

function RecommendationsSection({ data, loading, onUserAction }) {
  if (loading) {
    return (
      <LoadingMessage>
        <div>Generating personalized recommendations...</div>
      </LoadingMessage>
    );
  }

  if (!data || !data.recommendations) {
    return (
      <ErrorMessage>
        Failed to load recommendations. Please try again later.
      </ErrorMessage>
    );
  }

  const handleAction = (action, item) => {
    onUserAction(action, item.id, item.content_type);
  };

  const getContentTypeIcon = (type) => {
    switch (type) {
      case 'news': return '📰';
      case 'jobs': return '💼';
      case 'videos': return '🎥';
      case 'deals': return '🛍️';
      default: return '⭐';
    }
  };

  return (
    <RecommendationsContainer>
      {data.recommendations.slice(0, 5).map((item, index) => (
        <RecommendationItem key={item.id || index}>
          <RecommendationImage 
            src={item.image_url || item.thumbnail || '/placeholder-recommendation.jpg'} 
            alt={item.title}
            onError={(e) => {
              e.target.style.display = 'none';
            }}
          />
          <RecommendationContent>
            <RecommendationTitle>
              {getContentTypeIcon(item.content_type)} {item.title}
            </RecommendationTitle>
            <RecommendationDescription>
              {item.description?.substring(0, 100)}...
            </RecommendationDescription>
            <RecommendationMeta>
              <div>
                <span>{item.reason}</span>
                {item.recommendation_score && (
                  <RecommendationScore>
                    <FiStar />
                    {(item.recommendation_score * 100).toFixed(0)}%
                  </RecommendationScore>
                )}
              </div>
              <ActionButton 
                onClick={() => {
                  window.open(item.url, '_blank');
                  handleAction('click', item);
                }}
                title="View"
              >
                <FiExternalLink />
              </ActionButton>
            </RecommendationMeta>
          </RecommendationContent>
        </RecommendationItem>
      ))}
    </RecommendationsContainer>
  );
}

export default RecommendationsSection;
