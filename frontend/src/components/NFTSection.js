import React from 'react';
import styled from 'styled-components';
import { FiGift, FiClock, FiCheckCircle } from 'react-icons/fi';

const NFTsContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const NFTItem = styled.div`
  padding: 15px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  margin-bottom: 15px;
  background: rgba(255, 255, 255, 0.05);
`;

const NFTTitle = styled.h4`
  margin: 0 0 8px 0;
  font-size: 14px;
  color: white;
  display: flex;
  align-items: center;
  gap: 8px;
`;

const NFTDescription = styled.p`
  margin: 0 0 10px 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.4;
`;

const NFTMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
`;

const NFTStatus = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 10px;
  font-weight: 600;
  
  ${props => props.status === 'minted' && `
    background: rgba(34, 197, 94, 0.2);
    color: #22c55e;
    border: 1px solid rgba(34, 197, 94, 0.3);
  `}
  
  ${props => props.status === 'claimed' && `
    background: rgba(59, 130, 246, 0.2);
    color: #3b82f6;
    border: 1px solid rgba(59, 130, 246, 0.3);
  `}
`;

const DiscountBadge = styled.div`
  background: linear-gradient(45deg, #ff6b6b, #ffa500);
  color: white;
  padding: 4px 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
`;

const LoadingMessage = styled.div`
  text-align: center;
  padding: 40px;
  color: rgba(255, 255, 255, 0.6);
`;

const EmptyMessage = styled.div`
  text-align: center;
  padding: 40px;
  color: rgba(255, 255, 255, 0.6);
`;

function NFTSection({ data, loading, onUserAction }) {
  if (loading) {
    return <LoadingMessage>Loading your NFT rewards...</LoadingMessage>;
  }

  if (!data || !data.nfts || data.nfts.length === 0) {
    return (
      <EmptyMessage>
        <FiGift size={32} style={{ marginBottom: '10px', opacity: 0.5 }} />
        <div>No NFT rewards yet</div>
        <div style={{ fontSize: '12px', marginTop: '5px' }}>
          Engage with deals and jobs to earn NFT coupons!
        </div>
      </EmptyMessage>
    );
  }

  const handleAction = (action, nft) => {
    onUserAction(action, nft.id, 'nft');
  };

  return (
    <NFTsContainer>
      {data.nfts.map((nft, index) => (
        <NFTItem key={nft.id || index}>
          <NFTTitle>
            <FiGift />
            {nft.title}
          </NFTTitle>
          <NFTDescription>{nft.description}</NFTDescription>
          <NFTMeta>
            <div>
              <DiscountBadge>
                {nft.discount}% OFF
              </DiscountBadge>
              <div style={{ marginTop: '8px' }}>
                {nft.status === 'minted' && (
                  <>
                    <FiClock />
                    Expires: {new Date(nft.expires_at).toLocaleDateString()}
                  </>
                )}
                {nft.status === 'claimed' && (
                  <>
                    <FiCheckCircle />
                    Claimed: {new Date(nft.claimed_at).toLocaleDateString()}
                  </>
                )}
              </div>
            </div>
            <NFTStatus status={nft.status}>
              {nft.status === 'minted' && 'Available'}
              {nft.status === 'claimed' && 'Claimed'}
            </NFTStatus>
          </NFTMeta>
        </NFTItem>
      ))}
    </NFTsContainer>
  );
}

export default NFTSection;
