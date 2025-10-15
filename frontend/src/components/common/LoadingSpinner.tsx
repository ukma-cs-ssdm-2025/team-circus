import { Box, CircularProgress } from '@mui/material';

interface LoadingSpinnerProps {
  size?: number;
  py?: number;
}

const LoadingSpinner = ({ size, py = 4 }: LoadingSpinnerProps) => {
  return (
    <Box sx={{ display: 'flex', justifyContent: 'center', py }}>
      <CircularProgress size={size} />
    </Box>
  );
};

export default LoadingSpinner;
