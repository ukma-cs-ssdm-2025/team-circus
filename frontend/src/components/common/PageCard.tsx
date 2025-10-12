import { Card, CardContent, useTheme } from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface PageCardProps extends BaseComponentProps {
  children: React.ReactNode;
  padding?: number;
}

const PageCard = ({ children, className = '', padding = 6 }: PageCardProps) => {
  const theme = useTheme();

  return (
    <Card
      className={className}
      sx={{
        background: theme.palette.mode === 'light' 
          ? 'rgba(255, 255, 255, 0.8)' 
          : 'rgba(30, 30, 30, 0.8)',
        backdropFilter: 'blur(10px)',
        borderRadius: 4,
        boxShadow: '0 10px 30px rgba(0, 0, 0, 0.1)',
      }}
    >
      <CardContent sx={{ p: padding }}>
        {children}
      </CardContent>
    </Card>
  );
};

export default PageCard;
