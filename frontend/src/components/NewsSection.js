import React from 'react';
import styled from 'styled-components';
import { FiExternalLink, FiBookmark, FiShare2 } from 'react-icons/fi';

const NewsContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const NewsItem = styled.div`
  display: flex;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:last-child {
    border-bottom: none;
  }
`;

const NewsImage = styled.img`
  width: 80px;
  height: 60px;
  object-fit: cover;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
`;

const NewsContent = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const NewsTitle = styled.h4`
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

const NewsDescription = styled.p`
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const NewsMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
`;

const NewsActions = styled.div`
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

const ErrorMessage = styled.div`
  text-align: center;
  padding: 20px;
  color: #ff6b6b;
  background: rgba(255, 107, 107, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(255, 107, 107, 0.2);
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

function NewsSection({ data, loading, onUserAction }) {
  if (loading) {
    return (
      <LoadingMessage>
        <div>Loading latest news...</div>
      </LoadingMessage>
    );
  }

  if (!data || !data.articles) {
    return (
      <ErrorMessage>
        Failed to load news. Please try again later.
      </ErrorMessage>
    );
  }

  // Check if this is static data
  const isStaticData = data.error === "STATIC DATA - No API Key Found" || 
                      data.articles.some(article => article.is_static) ||
                      data.source !== "REAL API DATA";

  const handleAction = (action, article) => {
    onUserAction(action, article.id, 'news');
  };

  return (
    <NewsContainer>
      {isStaticData && (
        <StaticDataWarning>
          ⚠️ <strong>STATIC DATA WARNING:</strong> {data.message || 'This is mock data. Add NEWS_API_KEY to env.local for real news.'}
          <br />
          <small>Get API key: <a href="https://newsapi.org/register" target="_blank" rel="noopener noreferrer">https://newsapi.org/register</a></small>
        </StaticDataWarning>
      )}
      
      {!isStaticData && (
        <div style={{ marginBottom: '15px', textAlign: 'center' }}>
          <ApiStatusIndicator className="real-api">Real API Data</ApiStatusIndicator>
        </div>
      )}
      
      {data.articles.slice(0, 5).map((article, index) => (
        <NewsItem key={article.id || index}>
          <NewsImage 
            src={article.image_url || '/placeholder-news.jpg'} 
            alt={article.title}
            onError={(e) => {
              e.target.style.display = 'none';
            }}
          />
          <NewsContent>
            <NewsTitle>{article.title}</NewsTitle>
            <NewsDescription>
              {article.description?.substring(0, 100)}...
            </NewsDescription>
            <NewsMeta>
              <span>{article.source} • {new Date(article.published_at).toLocaleDateString()}</span>
              <NewsActions>
                <ActionButton 
                  onClick={() => handleAction('bookmark', article)}
                  title="Save Article"
                >
                  <FiBookmark />
                </ActionButton>
                <ActionButton 
                  onClick={() => handleAction('share', article)}
                  title="Share"
                >
                  <FiShare2 />
                </ActionButton>
                <ActionButton 
                  onClick={() => {
                    if (article.is_static) {
                      // For static data, redirect to API registration
                      window.open('https://newsapi.org/register', '_blank', 'noopener,noreferrer');
                    } else if (article.url && article.url !== 'https://example.com/news/1') {
                      // Real news article
                      window.open(article.url, '_blank', 'noopener,noreferrer');
                    } else {
                      // Fallback
                      window.open('https://news.google.com', '_blank', 'noopener,noreferrer');
                    }
                    handleAction('click', article);
                  }}
                  title={article.is_static ? "Get API Key" : "Read more"}
                >
                  <FiExternalLink />
                </ActionButton>
              </NewsActions>
            </NewsMeta>
          </NewsContent>
        </NewsItem>
      ))}
    </NewsContainer>
  );
}

export default NewsSection;