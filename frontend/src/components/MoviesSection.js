import React from 'react';
import styled from 'styled-components';
import { FiExternalLink, FiBookmark, FiShare2, FiStar } from 'react-icons/fi';

const MoviesContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const MovieItem = styled.div`
  display: flex;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:last-child {
    border-bottom: none;
  }
`;

const MoviePoster = styled.div`
  position: relative;
  width: 80px;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
`;

const PosterImage = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
`;

const RatingBadge = styled.div`
  position: absolute;
  top: 5px;
  right: 5px;
  background: rgba(0, 0, 0, 0.8);
  color: #ffc107;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 2px;
`;

const MovieContent = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const MovieTitle = styled.h4`
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

const MovieOverview = styled.p`
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const MovieMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
`;

const MovieActions = styled.div`
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

function MoviesSection({ data, loading, onUserAction }) {
  if (loading) {
    return <LoadingMessage>Loading latest movies...</LoadingMessage>;
  }

  if (!data || !data.movies) {
    return <LoadingMessage>No movies available at the moment.</LoadingMessage>;
  }

  // Check if this is static data
  const isStaticData = data.error === "STATIC DATA - No TMDB API Key Found" || 
                      data.movies.some(movie => movie.is_static) ||
                      data.source !== "REAL TMDB API DATA";

  const handleAction = (action, movie) => {
    onUserAction(action, movie.id, 'movies');
  };

  return (
    <MoviesContainer>
      {isStaticData && (
        <StaticDataWarning>
          ⚠️ <strong>STATIC DATA WARNING:</strong> {data.message || 'This is mock data. Add TMDB_API_KEY to env.local for real movies.'}
          <br />
          <small>Get API key: <a href="https://www.themoviedb.org/settings/api" target="_blank" rel="noopener noreferrer">https://www.themoviedb.org/settings/api</a></small>
        </StaticDataWarning>
      )}
      
      {!isStaticData && (
        <div style={{ marginBottom: '15px', textAlign: 'center' }}>
          <ApiStatusIndicator className="real-api">Real TMDB API</ApiStatusIndicator>
        </div>
      )}
      
      {data.movies.slice(0, 4).map((movie, index) => (
        <MovieItem key={movie.id || index}>
          <MoviePoster>
            <PosterImage 
              src={movie.poster_path || '/placeholder-movie.jpg'} 
              alt={movie.title}
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
            {movie.vote_average && (
              <RatingBadge>
                <FiStar size={8} />
                {movie.vote_average.toFixed(1)}
              </RatingBadge>
            )}
          </MoviePoster>
          <MovieContent>
            <MovieTitle>{movie.title}</MovieTitle>
            <MovieOverview>
              {movie.overview?.substring(0, 120)}...
            </MovieOverview>
            <MovieMeta>
              <span>{movie.release_date ? new Date(movie.release_date).getFullYear() : 'N/A'}</span>
              <MovieActions>
                <ActionButton 
                  onClick={() => handleAction('bookmark', movie)}
                  title="Save Movie"
                >
                  <FiBookmark />
                </ActionButton>
                <ActionButton 
                  onClick={() => handleAction('share', movie)}
                  title="Share"
                >
                  <FiShare2 />
                </ActionButton>
                <ActionButton 
                  onClick={() => {
                    if (movie.is_static) {
                      // For static data, redirect to API registration
                      window.open('https://www.themoviedb.org/settings/api', '_blank', 'noopener,noreferrer');
                    } else {
                      // Real movie - redirect to TMDB
                      window.open(`https://www.themoviedb.org/movie/${movie.id}`, '_blank', 'noopener,noreferrer');
                    }
                    handleAction('click', movie);
                  }}
                  title={movie.is_static ? "Get API Key" : "View Movie"}
                >
                  <FiExternalLink />
                </ActionButton>
              </MovieActions>
            </MovieMeta>
          </MovieContent>
        </MovieItem>
      ))}
    </MoviesContainer>
  );
}

export default MoviesSection;
