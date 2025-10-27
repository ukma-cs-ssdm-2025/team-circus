import { useCallback, useEffect, useState } from 'react';
import {
  addGroupMember,
  fetchGroupMembers,
  removeGroupMember,
  updateGroupMemberRole,
} from '../../groups/api';
import type { ApiError, GroupMember, GroupRole } from '../../../types';
import { normalizeApiError } from '../../../utils/apiError';
import { HttpError } from '../../../services/httpError';

interface UseGroupMembersResult {
  members: GroupMember[];
  loading: boolean;
  error: ApiError | null;
  mutating: boolean;
  refresh: () => Promise<void>;
  addMember: (userUUID: string, role: GroupRole) => Promise<void>;
  updateMemberRole: (userUUID: string, role: GroupRole) => Promise<void>;
  removeMember: (userUUID: string) => Promise<void>;
  reset: () => void;
}

export const useGroupMembers = (
  groupUUID: string | null,
): UseGroupMembersResult => {
  const [members, setMembers] = useState<GroupMember[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [mutating, setMutating] = useState(false);

  const loadMembers = useCallback(async () => {
    if (!groupUUID) {
      setMembers([]);
      setError(null);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const list = await fetchGroupMembers(groupUUID);
      setMembers(list);
    } catch (err) {
      setError(normalizeApiError(err));
    } finally {
      setLoading(false);
    }
  }, [groupUUID]);

  useEffect(() => {
    if (!groupUUID) {
      setMembers([]);
      setError(null);
      return;
    }

    void loadMembers();
  }, [groupUUID, loadMembers]);

  const wrapMutation = useCallback(
    async <T>(operation: () => Promise<T>): Promise<T> => {
      setMutating(true);
      setError(null);
      try {
        return await operation();
      } catch (err) {
        const apiError = normalizeApiError(err);
        setError(apiError);
        if (err instanceof Error) {
          throw err;
        }
        throw new HttpError(
          apiError.message,
          apiError.status,
          apiError.details,
          apiError.code,
        );
      } finally {
        setMutating(false);
      }
    },
    [],
  );

  const addMemberHandler = useCallback(
    async (userUUID: string, role: GroupRole) => {
      if (!groupUUID) {
        throw new Error('Group not selected');
      }

      await wrapMutation(() => addGroupMember(groupUUID, userUUID, role));
      await loadMembers();
    },
    [groupUUID, loadMembers, wrapMutation],
  );

  const updateMemberRoleHandler = useCallback(
    async (userUUID: string, role: GroupRole) => {
      if (!groupUUID) {
        throw new Error('Group not selected');
      }

      await wrapMutation(() =>
        updateGroupMemberRole(groupUUID, userUUID, role),
      );
      await loadMembers();
    },
    [groupUUID, loadMembers, wrapMutation],
  );

  const removeMemberHandler = useCallback(
    async (userUUID: string) => {
      if (!groupUUID) {
        throw new Error('Group not selected');
      }

      await wrapMutation(() => removeGroupMember(groupUUID, userUUID));
      await loadMembers();
    },
    [groupUUID, loadMembers, wrapMutation],
  );

  const reset = useCallback(() => {
    setMembers([]);
    setError(null);
    setLoading(false);
    setMutating(false);
  }, []);

  return {
    members,
    loading,
    error,
    mutating,
    refresh: loadMembers,
    addMember: addMemberHandler,
    updateMemberRole: updateMemberRoleHandler,
    removeMember: removeMemberHandler,
    reset,
  };
};
