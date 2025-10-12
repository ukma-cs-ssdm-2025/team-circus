import { Button as MuiButton, ButtonProps as MuiButtonProps } from '@mui/material';
import type { BaseComponentProps } from '../../types';

interface ButtonProps extends BaseComponentProps, Omit<MuiButtonProps, 'variant' | 'size'> {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'small' | 'medium' | 'large';
}

const Button = ({ 
  children, 
  className = '', 
  variant = 'primary', 
  size = 'medium',
  ...muiProps
}: ButtonProps) => {
  const getMuiVariant = () => {
    switch (variant) {
      case 'primary':
        return 'contained';
      case 'secondary':
        return 'outlined';
      case 'danger':
        return 'contained';
      default:
        return 'contained';
    }
  };

  const getMuiSize = () => {
    switch (size) {
      case 'small':
        return 'small';
      case 'medium':
        return 'medium';
      case 'large':
        return 'large';
      default:
        return 'medium';
    }
  };

  return (
    <MuiButton
      variant={getMuiVariant()}
      size={getMuiSize()}
      className={className}
      color={variant === 'danger' ? 'error' : 'primary'}
      {...muiProps}
    >
      {children}
    </MuiButton>
  );
};

export default Button;
