import { Box, Typography, useTheme } from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface PageHeaderProps extends BaseComponentProps {
  title: string;
  subtitle?: string;
  variant?: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6';
}

const PageHeader = ({ 
  title, 
  subtitle, 
  variant = 'h3', 
  className = '' 
}: PageHeaderProps) => {
  const theme = useTheme();

  return (
    <Box className={className}>
      <Typography
        variant={variant}
        sx={{
          fontWeight: 700,
          color: theme.palette.primary.main,
          mb: subtitle ? 1 : 0,
        }}
      >
        {title}
      </Typography>
      
      {subtitle && (
        <Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
          {subtitle}
        </Typography>
      )}
    </Box>
  );
};

export default PageHeader;
