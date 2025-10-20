import { Box, Stack } from '@mui/material';
import { ActionButton } from '../components/forms';
import { useLanguage } from '../contexts/LanguageContext';
import type { BaseComponentProps } from '../types';

type HomeProps = BaseComponentProps;

const Home = ({ className = '' }: HomeProps) => {
  const { t } = useLanguage();

  return (
    <Box
      className={className}
      sx={{
        flex: 1,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: { xs: '60vh', md: '100%' },
        px: { xs: 2, md: 4 },
        py: { xs: 3, md: 4 },
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
