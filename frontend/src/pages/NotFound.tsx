import { Link } from 'react-router-dom';
import {
  Box,
  Container,
  Typography,
  Button,
  Stack,
  useTheme
} from '@mui/material';
import { Home as HomeIcon, ArrowBack as ArrowBackIcon } from '@mui/icons-material';
import { useLanguage } from '../contexts/LanguageContext';
import { ROUTES } from '../constants';
import type { BaseComponentProps } from '../types';

interface NotFoundProps extends BaseComponentProps {}

const NotFound = ({ className = '' }: NotFoundProps) => {
  const theme = useTheme();
  const { t } = useLanguage();

  return (
    <Box className={className}>
      <Container maxWidth="sm" sx={{ py: 8 }}>
        <Box
          sx={{
            textAlign: 'center',
            background: theme.palette.mode === 'light' 
              ? 'rgba(255, 255, 255, 0.8)' 
              : 'rgba(30, 30, 30, 0.8)',
            backdropFilter: 'blur(10px)',
            borderRadius: 4,
            p: 6,
            boxShadow: '0 10px 30px rgba(0, 0, 0, 0.1)',
          }}
        >
          <Typography
            variant="h1"
            sx={{
              fontSize: '8rem',
              fontWeight: 900,
              color: theme.palette.primary.main,
              mb: 2,
              lineHeight: 1,
            }}
          >
            404
          </Typography>
          
          <Typography variant="h3" gutterBottom>
            {t('notFound.title')}
          </Typography>
          
          <Typography variant="h6" color="text.secondary" sx={{ mb: 4 }}>
            {t('notFound.message')}
          </Typography>

          <Stack
            direction={{ xs: 'column', sm: 'row' }}
            spacing={2}
            justifyContent="center"
          >
            <Button
              component={Link}
              to={ROUTES.HOME}
              variant="contained"
              size="large"
              startIcon={<HomeIcon />}
              sx={{ minWidth: 200 }}
            >
              {t('notFound.home')}
            </Button>
            <Button
              variant="outlined"
              size="large"
              startIcon={<ArrowBackIcon />}
              onClick={() => window.history.back()}
              sx={{ minWidth: 200 }}
            >
              {t('notFound.back')}
            </Button>
          </Stack>
        </Box>
      </Container>
    </Box>
  );
};

export default NotFound;
