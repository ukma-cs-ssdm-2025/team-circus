import { 
  Box, 
  Button, 
  Stack,
  useTheme
} from '@mui/material';
import { useLanguage } from '../contexts/LanguageContext';
import type { BaseComponentProps } from '../types';

interface HomeProps extends BaseComponentProps {}

const Home = ({ className = '' }: HomeProps) => {
  const theme = useTheme();
  const { t } = useLanguage();

  return (
    <Box 
      className={className}
      sx={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '100vh',
        padding: 2,
      }}
    >
      <Stack
        direction={{ xs: 'column', sm: 'row' }}
        spacing={3}
        justifyContent="center"
        alignItems="center"
      >
        <Button
          variant="contained"
          size="large"
          sx={{ 
            minWidth: 200,
            height: 60,
            fontSize: '1.1rem',
            backgroundColor: theme.palette.primary.main,
            '&:hover': {
              backgroundColor: theme.palette.primary.dark,
            }
          }}
          onClick={() => alert(t('home.createDocumentAlert'))}
        >
          {t('home.createDocument')}
        </Button>
        <Button
          variant="contained"
          size="large"
          sx={{ 
            minWidth: 200,
            height: 60,
            fontSize: '1.1rem',
            backgroundColor: theme.palette.primary.main,
            '&:hover': {
              backgroundColor: theme.palette.primary.dark,
            }
          }}
          onClick={() => alert(t('home.createGroupAlert'))}
        >
          {t('home.createGroup')}
        </Button>
      </Stack>
    </Box>
  );
};

export default Home;
