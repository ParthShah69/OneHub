import React from 'react';
import styled from 'styled-components';
import { FiExternalLink, FiBookmark, FiShare2 } from 'react-icons/fi';

const VideosContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const VideoItem = styled.div`
  display: flex;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:last-child {
    border-bottom: none;
  }
`;

const VideoThumbnail = styled.div`
  position: relative;
  width: 120px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
`;

const ThumbnailImage = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
`;

const PlayButton = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 30px;
  height: 30px;
  background: rgba(0, 0, 0, 0.7);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
`;

const VideoContent = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const VideoTitle = styled.h4`
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

const VideoDescription = styled.p`
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const VideoMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
`;

const VideoActions = styled.div`
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

function VideosSection({ data, loading, onUserAction }) {
  if (loading) {
    return <LoadingMessage>Loading trending videos...</LoadingMessage>;
  }

  if (!data || !data.videos) {
    return <LoadingMessage>No videos available at the moment.</LoadingMessage>;
  }

  // Check if this is static data
  const isStaticData = data.error === "STATIC DATA - No YouTube API Key Found" || 
                      data.videos.some(video => video.is_static) ||
                      data.source !== "REAL YOUTUBE API DATA";

  const handleAction = (action, video) => {
    onUserAction(action, video.id, 'videos');
  };

  return (
    <VideosContainer>
      {isStaticData && (
        <StaticDataWarning>
          ⚠️ <strong>STATIC DATA WARNING:</strong> {data.message || 'This is mock data. Add YOUTUBE_API_KEY to env.local for real videos.'}
          <br />
          <small>Get API key: <a href="https://console.developers.google.com/" target="_blank" rel="noopener noreferrer">https://console.developers.google.com/</a></small>
        </StaticDataWarning>
      )}
      
      {!isStaticData && (
        <div style={{ marginBottom: '15px', textAlign: 'center' }}>
          <ApiStatusIndicator className="real-api">Real YouTube API</ApiStatusIndicator>
        </div>
      )}
      
      {data.videos.slice(0, 5).map((video, index) => (
        <VideoItem key={video.id || index}>
          <VideoThumbnail>
            <ThumbnailImage 
              src={video.thumbnail || '/placeholder-video.jpg'} 
              alt={video.title}
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
            <PlayButton>▶</PlayButton>
          </VideoThumbnail>
          <VideoContent>
            <VideoTitle>{video.title}</VideoTitle>
            <VideoDescription>
              {video.description?.substring(0, 100)}...
            </VideoDescription>
            <VideoMeta>
              <span>{video.channel} • {video.views?.toLocaleString()} views</span>
              <VideoActions>
                <ActionButton 
                  onClick={() => handleAction('bookmark', video)}
                  title="Save Video"
                >
                  <FiBookmark />
                </ActionButton>
                <ActionButton 
                  onClick={() => handleAction('share', video)}
                  title="Share"
                >
                  <FiShare2 />
                </ActionButton>
                <ActionButton 
                  onClick={() => {
                    if (video.is_static) {
                      // For static data, redirect to API registration
                      window.open('https://console.developers.google.com/', '_blank', 'noopener,noreferrer');
                    } else if (video.url && video.url !== 'https://youtube.com/watch?v=demo1') {
                      // Real YouTube video
                      window.open(video.url, '_blank', 'noopener,noreferrer');
                    } else {
                      // Fallback
                      window.open('https://youtube.com', '_blank', 'noopener,noreferrer');
                    }
                    handleAction('click', video);
                  }}
                  title={video.is_static ? "Get API Key" : "Watch"}
                >
                  <FiExternalLink />
                </ActionButton>
              </VideoActions>
            </VideoMeta>
          </VideoContent>
        </VideoItem>
      ))}
    </VideosContainer>
  );
}

export default VideosSection;