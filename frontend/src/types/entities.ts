// Entity types for the application

export interface GroupItem {
  uuid: string;
  name: string;
  created_at: string;
}

export interface DocumentItem {
  uuid: string;
  name: string;
  content: string;
  group_uuid: string;
  created_at: string;
}

export interface GroupsResponse {
  groups: GroupItem[];
}

export interface DocumentsResponse {
  documents: DocumentItem[];
}

// Filter and search types
export interface GroupOption {
  value: string;
  label: string;
}

export interface DocumentFilters {
  selectedGroup: string;
  searchTerm: string;
}
