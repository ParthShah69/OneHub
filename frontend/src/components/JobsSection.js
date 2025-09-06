import React from 'react';
import styled from 'styled-components';
import { FiExternalLink, FiBookmark, FiShare2 } from 'react-icons/fi';

const JobsContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const JobItem = styled.div`
  padding: 15px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  margin-bottom: 15px;
  background: rgba(255, 255, 255, 0.02);
  transition: all 0.2s ease;

  &:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.2);
  }
`;

const JobHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 10px;
`;

const JobTitle = styled.h4`
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: white;
  line-height: 1.3;
`;

const JobCompany = styled.p`
  margin: 5px 0 0 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
`;

const JobLocation = styled.span`
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  background: rgba(255, 255, 255, 0.1);
  padding: 4px 8px;
  border-radius: 12px;
`;

const JobDescription = styled.p`
  margin: 10px 0;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const JobMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 10px;
`;

const JobActions = styled.div`
  display: flex;
  gap: 8px;
`;

const ActionButton = styled.button`
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  padding: 6px;
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

function JobsSection({ data, loading, onUserAction }) {
  if (loading) {
    return <LoadingMessage>Loading job opportunities...</LoadingMessage>;
  }

  if (!data || !data.jobs) {
    return <LoadingMessage>No jobs available at the moment.</LoadingMessage>;
  }

  // Check if this is static data
  const isStaticData = data.error === "STATIC DATA - No LinkedIn API Key Found" || 
                      data.jobs.some(job => job.is_static) ||
                      data.source !== "REAL LINKEDIN API DATA";

  const handleAction = (action, job) => {
    onUserAction(action, job.id, 'jobs');
  };

  return (
    <JobsContainer>
      {isStaticData && (
        <StaticDataWarning>
          ⚠️ <strong>STATIC DATA WARNING:</strong> {data.message || 'This is mock data. Add LINKEDIN_API_KEY to env.local for real jobs.'}
          <br />
          <small>Get API key: <a href="https://developer.linkedin.com/" target="_blank" rel="noopener noreferrer">https://developer.linkedin.com/</a></small>
        </StaticDataWarning>
      )}
      
      {!isStaticData && (
        <div style={{ marginBottom: '15px', textAlign: 'center' }}>
          <ApiStatusIndicator className="real-api">Real LinkedIn API</ApiStatusIndicator>
        </div>
      )}
      
      {data.jobs.slice(0, 3).map((job, index) => (
        <JobItem key={job.id || index}>
          <JobHeader>
            <div>
              <JobTitle>{job.title}</JobTitle>
              <JobCompany>{job.company}</JobCompany>
            </div>
            <JobLocation>{job.location}</JobLocation>
          </JobHeader>
          
          <JobDescription>
            {job.description?.substring(0, 150)}...
          </JobDescription>
          
          <JobMeta>
            <span>{job.type} • {job.salary}</span>
            <JobActions>
              <ActionButton 
                onClick={() => handleAction('bookmark', job)}
                title="Save Job"
              >
                <FiBookmark />
              </ActionButton>
              <ActionButton 
                onClick={() => handleAction('share', job)}
                title="Share"
              >
                <FiShare2 />
              </ActionButton>
              <ActionButton 
                onClick={() => {
                  if (job.is_static) {
                    // For static data, redirect to API registration
                    window.open('https://developer.linkedin.com/', '_blank', 'noopener,noreferrer');
                  } else if (job.url && job.url !== 'https://example.com/jobs/1') {
                    // Real job posting
                    window.open(job.url, '_blank', 'noopener,noreferrer');
                  } else {
                    // Fallback
                    window.open('https://linkedin.com/jobs', '_blank', 'noopener,noreferrer');
                  }
                  handleAction('click', job);
                }}
                title={job.is_static ? "Get API Key" : "Apply"}
              >
                <FiExternalLink />
              </ActionButton>
            </JobActions>
          </JobMeta>
        </JobItem>
      ))}
    </JobsContainer>
  );
}

export default JobsSection;