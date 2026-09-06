import React, { useCallback, useEffect } from 'react';

export const DataSyncManager = ({ id }: { id: string }) => {
  const handleSync = useCallback(() => {
    console.log('syncing', id);
  }, [id]);

  useEffect(() => {
    handleSync();
  }, [handleSync]);

  return <div>Sync active for {id}</div>;
};
