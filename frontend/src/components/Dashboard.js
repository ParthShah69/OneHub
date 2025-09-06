import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import axios from 'axios';
import NewsSection from './NewsSection';
import JobsSection from './JobsSection';
import VideosSection from './VideosSection';
import DealsSection from './DealsSection';
import MoviesSection from './MoviesSection';
import FoodSection from './FoodSection';
import RecommendationsSection from './RecommendationsSection';
import NFTSection from './NFTSection';
import { FiRefreshCw, FiSettings } from 'react-icons/fi';

const DashboardContainer = styled.div`
  color: white;
`;

const DashboardHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  backdrop-filter: blur(10px);
`;

const WelcomeText = styled.h1`
  margin: 0;
  font-size: 2.5rem;
  font-weight: 300;
`;

const ActionButtons = styled.div`
  display: flex;
  gap: 15px;
`;

const RefreshButton = styled.button`
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 25px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 14px;

  &:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: translateY(-2px);
  }
`;

const SettingsButton = styled(RefreshButton)`
  background: rgba(255, 255, 255, 0.15);
`;

const GridContainer = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 25px;
  margin-bottom: 30px;
`;

const SectionCard = styled.div`
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 25px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  }
`;

const SectionTitle = styled.h2`
  margin: 0 0 20px 0;
  font-size: 1.5rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 10px;
`;

const LoadingSpinner = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  font-size: 18px;
`;

const ErrorMessage = styled.div`
  background: rgba(255, 0, 0, 0.2);
  border: 1px solid rgba(255, 0, 0, 0.3);
  border-radius: 10px;
  padding: 15px;
  margin: 10px 0;
  color: #ff6b6b;
`;

function Dashboard({ user }) {
  const [data, setData] = useState({
    news: null,
    jobs: null,
    videos: null,
    deals: null,
    movies: null,
    food: null,
    recommendations: null,
    nfts: null
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [lastRefresh, setLastRefresh] = useState(new Date());

  const fetchData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const [newsRes, jobsRes, videosRes, dealsRes, moviesRes, foodRes, recommendationsRes, nftsRes] = await Promise.allSettled([
        axios.get(`/api/news/trending?user_id=${user.id}`),
        axios.get(`/api/jobs/trending?user_id=${user.id}`),
        axios.get(`/api/videos/trending?user_id=${user.id}`),
        axios.get(`/api/deals/trending?user_id=${user.id}`),
        axios.get(`/api/movies/trending?user_id=${user.id}`),
        axios.get(`/api/food/trending?user_id=${user.id}`),
        axios.get(`/api/recommendations?user_id=${user.id}`),
        axios.get(`/api/nft/${user.id}`)
      ]);

      setData({
        news: newsRes.status === 'fulfilled' ? newsRes.value.data : null,
        jobs: jobsRes.status === 'fulfilled' ? jobsRes.value.data : null,
        videos: videosRes.status === 'fulfilled' ? videosRes.value.data : null,
        deals: dealsRes.status === 'fulfilled' ? dealsRes.value.data : null,
        movies: moviesRes.status === 'fulfilled' ? moviesRes.value.data : null,
        food: foodRes.status === 'fulfilled' ? foodRes.value.data : null,
        recommendations: recommendationsRes.status === 'fulfilled' ? recommendationsRes.value.data : null,
        nfts: nftsRes.status === 'fulfilled' ? nftsRes.value.data : null
      });

      setLastRefresh(new Date());
    } catch (err) {
      setError('Failed to fetch data. Please try again.');
      console.error('Error fetching data:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [user]);

  const handleRefresh = () => {
    fetchData();
  };

  const handleUserAction = (action, contentId, category) => {
    // Track user behavior
    axios.post(`/api/users/${user.id}/behavior`, {
      action,
      content_id: contentId,
      category
    }).catch(err => console.error('Failed to track behavior:', err));
  };

  if (loading && !data.news) {
    return (
      <DashboardContainer>
        <LoadingSpinner>
          <FiRefreshCw className="spinning" />
          Loading your personalized dashboard...
        </LoadingSpinner>
      </DashboardContainer>
    );
  }

  return (
    <DashboardContainer>
      <DashboardHeader>
        <WelcomeText>
          Welcome back, {user.name}! 👋
        </WelcomeText>
        <ActionButtons>
          <RefreshButton onClick={handleRefresh}>
            <FiRefreshCw />
            Refresh
          </RefreshButton>
          <SettingsButton>
            <FiSettings />
            Settings
          </SettingsButton>
        </ActionButtons>
      </DashboardHeader>

      {error && <ErrorMessage>{error}</ErrorMessage>}

      <GridContainer>
        <SectionCard>
          <SectionTitle>🎯 Personalized Recommendations</SectionTitle>
          <RecommendationsSection 
            data={data.recommendations} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>📰 Latest News</SectionTitle>
          <NewsSection 
            data={data.news} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>💼 Job Opportunities</SectionTitle>
          <JobsSection 
            data={data.jobs} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>🎥 Trending Videos</SectionTitle>
          <VideosSection 
            data={data.videos} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>🛍️ Hot Deals</SectionTitle>
          <DealsSection 
            data={data.deals} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>🎬 Latest Movies</SectionTitle>
          <MoviesSection 
            data={data.movies} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>🍽️ Delicious Recipes</SectionTitle>
          <FoodSection 
            data={data.food} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>

        <SectionCard>
          <SectionTitle>🎁 NFT Rewards</SectionTitle>
          <NFTSection 
            data={data.nfts} 
            loading={loading}
            onUserAction={handleUserAction}
          />
        </SectionCard>
      </GridContainer>

      <div style={{ textAlign: 'center', color: 'rgba(255, 255, 255, 0.7)', fontSize: '14px' }}>
        Last updated: {lastRefresh.toLocaleTimeString()}
      </div>
    </DashboardContainer>
  );
}

export default Dashboard;
