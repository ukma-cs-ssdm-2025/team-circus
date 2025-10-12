import { Box, Stack } from '@mui/material';
import { ActionButton } from '../components/forms';
import { useLanguage } from '../contexts/LanguageContext';
import type { BaseComponentProps } from '../types';

interface HomeProps extends BaseComponentProps {}

const Home = ({ className = '' }: HomeProps) => {
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
        <ActionButton
          onClick={() => alert(t('home.createDocumentAlert'))}
        >
          {t('home.createDocument')}
        </ActionButton>
        <ActionButton
          onClick={() => alert(t('home.createGroupAlert'))}
        >
          {t('home.createGroup')}
        </ActionButton>
      </Stack>
    </Box>
  );
};

export default Home;
