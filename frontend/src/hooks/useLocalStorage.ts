import { useState } from 'react';
import { storage } from '../utils';

export function useLocalStorage<T>(key: string, initialValue: T) {
  const [storedValue, setStoredValue] = useState<T>(() => {
    return storage.get(key, initialValue);
  });

  const setValue = (value: T | ((val: T) => T)) => {
    try {
      let valueToStore: T;
      if (typeof value === 'function') {
        const updater = value as (val: T) => T;
        valueToStore = updater(storedValue);
      } else {
        valueToStore = value;
      }
      setStoredValue(valueToStore);
      storage.set(key, valueToStore);
    } catch (error) {
      console.error('Error setting localStorage:', error);
    }
  };

  return [storedValue, setValue] as const;
}
