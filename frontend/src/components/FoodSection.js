import React from 'react';
import styled from 'styled-components';
import { FiExternalLink, FiBookmark, FiShare2, FiClock, FiUsers, FiHeart } from 'react-icons/fi';

const FoodContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
`;

const RecipeItem = styled.div`
  display: flex;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  
  &:last-child {
    border-bottom: none;
  }
`;

const RecipeImage = styled.div`
  position: relative;
  width: 100px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
`;

const Image = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
`;

const HealthBadge = styled.div`
  position: absolute;
  top: 5px;
  right: 5px;
  background: rgba(76, 175, 80, 0.9);
  color: white;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 2px;
`;

const RecipeContent = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const RecipeTitle = styled.h4`
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

const RecipeSummary = styled.p`
  margin: 0;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const RecipeMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
`;

const RecipeInfo = styled.div`
  display: flex;
  gap: 12px;
  align-items: center;
`;

const InfoItem = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 10px;
`;

const RecipeActions = styled.div`
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

function FoodSection({ data, loading, onUserAction }) {
  if (loading) {
    return <LoadingMessage>Loading delicious recipes...</LoadingMessage>;
  }

  if (!data || !data.recipes) {
    return <LoadingMessage>No recipes available at the moment.</LoadingMessage>;
  }

  // Check if this is static data
  const isStaticData = data.error === "STATIC DATA - No Recipe API Key Found" || 
                      data.recipes.some(recipe => recipe.is_static) ||
                      data.source !== "REAL SPOONACULAR API DATA";

  const handleAction = (action, recipe) => {
    onUserAction(action, recipe.id, 'food');
  };

  return (
    <FoodContainer>
      {isStaticData && (
        <StaticDataWarning>
          ⚠️ <strong>STATIC DATA WARNING:</strong> {data.message || 'This is mock data. Add RECIPE_API_KEY to env.local for real recipes.'}
          <br />
          <small>Get API key: <a href="https://spoonacular.com/food-api" target="_blank" rel="noopener noreferrer">https://spoonacular.com/food-api</a></small>
        </StaticDataWarning>
      )}
      
      {!isStaticData && (
        <div style={{ marginBottom: '15px', textAlign: 'center' }}>
          <ApiStatusIndicator className="real-api">Real Spoonacular API</ApiStatusIndicator>
        </div>
      )}
      
      {data.recipes.slice(0, 4).map((recipe, index) => (
        <RecipeItem key={recipe.id || index}>
          <RecipeImage>
            <Image 
              src={recipe.image || '/placeholder-recipe.jpg'} 
              alt={recipe.title}
              onError={(e) => {
                e.target.style.display = 'none';
              }}
            />
            {recipe.healthScore && (
              <HealthBadge>
                <FiHeart size={8} />
                {recipe.healthScore}
              </HealthBadge>
            )}
          </RecipeImage>
          <RecipeContent>
            <RecipeTitle>{recipe.title}</RecipeTitle>
            <RecipeSummary>
              {recipe.summary?.replace(/<[^>]*>/g, '').substring(0, 100)}...
            </RecipeSummary>
            <RecipeMeta>
              <RecipeInfo>
                {recipe.readyInMinutes && (
                  <InfoItem>
                    <FiClock size={10} />
                    {recipe.readyInMinutes}m
                  </InfoItem>
                )}
                {recipe.servings && (
                  <InfoItem>
                    <FiUsers size={10} />
                    {recipe.servings}
                  </InfoItem>
                )}
              </RecipeInfo>
              <RecipeActions>
                <ActionButton 
                  onClick={() => handleAction('bookmark', recipe)}
                  title="Save Recipe"
                >
                  <FiBookmark />
                </ActionButton>
                <ActionButton 
                  onClick={() => handleAction('share', recipe)}
                  title="Share"
                >
                  <FiShare2 />
                </ActionButton>
                <ActionButton 
                  onClick={() => {
                    if (recipe.is_static) {
                      // For static data, redirect to API registration
                      window.open('https://spoonacular.com/food-api', '_blank', 'noopener,noreferrer');
                    } else {
                      // Real recipe - redirect to Spoonacular
                      window.open(`https://spoonacular.com/recipe/${recipe.id}`, '_blank', 'noopener,noreferrer');
                    }
                    handleAction('click', recipe);
                  }}
                  title={recipe.is_static ? "Get API Key" : "View Recipe"}
                >
                  <FiExternalLink />
                </ActionButton>
              </RecipeActions>
            </RecipeMeta>
          </RecipeContent>
        </RecipeItem>
      ))}
    </FoodContainer>
  );
}

export default FoodSection;
