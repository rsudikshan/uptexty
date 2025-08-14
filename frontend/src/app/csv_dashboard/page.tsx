"use client";
import { useState, useEffect, ChangeEvent, useCallback, useMemo, memo } from "react";
import {
  Box,
  Button,
  TextField,
  Typography,
  List,
  ListItem,
  ListItemText,
  Paper,
  CircularProgress,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Chip,
  Pagination,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from "@mui/material";
import {
  ArrowBack,
  Add,
  Delete,
  DragIndicator,
  KeyboardArrowUp,
  KeyboardArrowDown,
} from "@mui/icons-material";

interface FileItem {
  id: number;
  filename: string;
  uploaded_at: string;
}

interface CSVRow {
  id: number;
  position: number;
  input_text: string;
}

// Memoized row component to prevent unnecessary re-renders
const CSVRowComponent = memo(({ 
  row, 
  index, 
  globalIndex,
  draggedRowIndex,
  dragOverIndex,
  onDragStart,
  onDragOver,
  onDragLeave,
  onDrop,
  onInsertRow,
  onDeleteRow,
  onEditRow 
}: {
  row: CSVRow;
  index: number;
  globalIndex: number;
  draggedRowIndex: number | null;
  dragOverIndex: number | null;
  onDragStart: (e: React.DragEvent, index: number) => void;
  onDragOver: (e: React.DragEvent, index: number) => void;
  onDragLeave: () => void;
  onDrop: (e: React.DragEvent, dropIndex: number) => void;
  onInsertRow: (index: number) => void;
  onDeleteRow: (row: CSVRow, index: number) => void;
  onEditRow: (row: CSVRow, newText: string) => void;
}) => {
  const [localText, setLocalText] = useState(row.input_text);
  const [hasChanges, setHasChanges] = useState(false);

  // Update local text when row changes
  useEffect(() => {
    setLocalText(row.input_text);
    setHasChanges(false);
  }, [row.input_text]);

  const handleTextChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setLocalText(newValue);
    setHasChanges(newValue !== row.input_text);
  }, [row.input_text]);

  const handleBlur = useCallback(() => {
    if (hasChanges && localText !== row.input_text) {
      onEditRow(row, localText);
      setHasChanges(false);
    }
  }, [hasChanges, localText, row, onEditRow]);

  return (
    <TableRow
      draggable
      onDragStart={(e) => onDragStart(e, globalIndex)}
      onDragOver={(e) => onDragOver(e, globalIndex)}
      onDragLeave={onDragLeave}
      onDrop={(e) => onDrop(e, globalIndex)}
      sx={{
        cursor: "move",
        "&:hover": { backgroundColor: "#f5f5f5" },
        backgroundColor: 
          draggedRowIndex === globalIndex ? "#e3f2fd" : 
          dragOverIndex === globalIndex ? "#f3e5f5" : "inherit",
        borderLeft: dragOverIndex === globalIndex ? "3px solid #9c27b0" : "none",
      }}
    >
      <TableCell sx={{ width: 120 }}>
        <Box display="flex" gap={1}>
          <IconButton size="small" sx={{ cursor: "grab" }}>
            <DragIndicator />
          </IconButton>
          <IconButton 
            size="small" 
            color="primary"
            onClick={() => onInsertRow(globalIndex)}
            title="Insert row before this one"
          >
            <Add />
          </IconButton>
          <IconButton 
            size="small" 
            color="error"
            onClick={() => onDeleteRow(row, globalIndex)}
            title="Delete this row"
          >
            <Delete />
          </IconButton>
        </Box>
      </TableCell>
      <TableCell sx={{ width: 80 }}>
        <Typography variant="body2" color="textSecondary">
          {globalIndex + 1}
        </Typography>
      </TableCell>
      <TableCell sx={{ width: 100 }}>
        <Typography variant="body2" color="textSecondary">
          {row.position.toFixed(3)}
        </Typography>
      </TableCell>
      <TableCell>
        <TextField
          fullWidth
          multiline
          variant="outlined"
          value={localText}
          onChange={handleTextChange}
          onBlur={handleBlur}
          size="small"
          sx={{ 
            minWidth: 300,
            '& .MuiInputBase-root': {
              backgroundColor: hasChanges ? '#fff3e0' : 'inherit'
            }
          }}
        />
      </TableCell>
    </TableRow>
  );
});

CSVRowComponent.displayName = 'CSVRowComponent';

export default function FilesPage() {
  const [files, setFiles] = useState<FileItem[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [filename, setFilename] = useState("");
  const [loading, setLoading] = useState(false);
  
  // CSV data display state
  const [currentFile, setCurrentFile] = useState<FileItem | null>(null);
  const [allCsvRows, setAllCsvRows] = useState<CSVRow[]>([]); // All rows
  const [totalRows, setTotalRows] = useState(0);
  const [currentPage, setCurrentPage] = useState(1);
  const [rowsPerPage, setRowsPerPage] = useState(100);
  const [insertDialogOpen, setInsertDialogOpen] = useState(false);
  const [insertPosition, setInsertPosition] = useState(0);
  const [newRowText, setNewRowText] = useState("");
  const [draggedRowIndex, setDraggedRowIndex] = useState<number | null>(null);
  const [dragOverIndex, setDragOverIndex] = useState<number | null>(null);

  const token = localStorage.getItem("jwt");

  // Calculate pagination
  const totalPages = Math.ceil(totalRows / rowsPerPage);
  const startIndex = (currentPage - 1) * rowsPerPage;
  const endIndex = startIndex + rowsPerPage;

  // Get current page rows
  const currentPageRows = useMemo(() => {
    return allCsvRows.slice(startIndex, endIndex);
  }, [allCsvRows, startIndex, endIndex]);

  // Fetch uploaded files
  const fetchFiles = useCallback(async () => {
    if (!token) return;
    setLoading(true);
    try {
      const res = await fetch("http://localhost:8080/files", {
        headers: {
          Authorization: "Bearer " + token,
        },
      });
      const data = await res.json();
      setFiles(data.body || []);
    } catch (err) {
      console.error("Error fetching files:", err);
    } finally {
      setLoading(false);
    }
  }, [token]);

  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  // Handle file input
  const handleFileChange = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  }, []);

  // Upload CSV file
  const handleUpload = useCallback(async () => {
    if (!selectedFile || !filename) {
      alert("Please select a file and enter a filename.");
      return;
    }

    const formData = new FormData();
    formData.append("file", selectedFile);
    formData.append("filename", filename);

    try {
      const res = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
        headers: {
          Authorization: "Bearer " + token,
        },
      });

      const data = await res.json();
      if (data.status) {
        setSelectedFile(null);
        setFilename("");
        fetchFiles();
      } else {
        alert("Upload failed: " + data.message);
      }
    } catch (err) {
      console.error("Upload error:", err);
      alert("Upload failed");
    }
  }, [selectedFile, filename, token, fetchFiles]);

  // Fetch CSV rows for display
  const fetchCSVRows = useCallback(async (file: FileItem, maintainPage = false) => {
    if (!token) return;
    setLoading(true);
    try {
      const res = await fetch(`http://localhost:8080/files/${file.id}`, {
        headers: {
          Authorization: "Bearer " + token,
        },
      });
      const data = await res.json();
      
      if (data.status) {
        const sortedRows = (data.body || []).sort((a: CSVRow, b: CSVRow) => a.position - b.position);
        setAllCsvRows(sortedRows);
        setTotalRows(sortedRows.length);
        setCurrentFile(file);
        
        // Only reset to first page if not maintaining current page
        if (!maintainPage) {
          setCurrentPage(1);
        }
      } else {
        alert("Failed to load CSV data: " + data.message);
      }
    } catch (err) {
      console.error("Error fetching CSV data:", err);
      alert("Failed to load CSV data");
    } finally {
      setLoading(false);
    }
  }, [token]);

  // Handle file click to display CSV data
  const handleFileClick = useCallback((file: FileItem) => {
    fetchCSVRows(file);
  }, [fetchCSVRows]);

  // Handle back to file list
  const handleBackToFiles = useCallback(() => {
    setCurrentFile(null);
    setAllCsvRows([]);
    setTotalRows(0);
    setCurrentPage(1);
  }, []);

  // Calculate position for new row insertion
  const calculateInsertPosition = useCallback((globalIndex: number): number => {
    if (allCsvRows.length === 0) return 1.0;
    
    if (globalIndex === 0) {
      return allCsvRows[0].position / 2;
    } else if (globalIndex >= allCsvRows.length) {
      return allCsvRows[allCsvRows.length - 1].position + 1.0;
    } else {
      const prevPosition = allCsvRows[globalIndex - 1].position;
      const nextPosition = allCsvRows[globalIndex].position;
      return (prevPosition + nextPosition) / 2;
    }
  }, [allCsvRows]);

  // Handle insert row dialog
  const handleInsertRow = useCallback((globalIndex: number) => {
    setInsertPosition(globalIndex);
    setNewRowText("");
    setInsertDialogOpen(true);
  }, []);

  // Confirm insert row
  const confirmInsertRow = useCallback(async () => {
    if (!token || !currentFile || !newRowText.trim()) return;
    
    const position = calculateInsertPosition(insertPosition);
    
    try {
      const res = await fetch(`http://localhost:8080/files/${currentFile.id}/rows`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify({
          position: position,
          input_text: newRowText.trim(),
        }),
      });

      const data = await res.json();
      if (data.status) {
        fetchCSVRows(currentFile, true); // Maintain current page
        setInsertDialogOpen(false);
      } else {
        alert("Failed to insert row: " + data.message);
      }
    } catch (err) {
      console.error("Error inserting row:", err);
      alert("Failed to insert row");
    }
  }, [token, currentFile, newRowText, insertPosition, calculateInsertPosition, fetchCSVRows]);

  // Handle page change
  const handlePageChange = useCallback((_: React.ChangeEvent<unknown>, page: number) => {
    setCurrentPage(page);
  }, []);

  // Handle rows per page change
  const handleRowsPerPageChange = useCallback((event: any) => {
    const newRowsPerPage = event.target.value;
    setRowsPerPage(newRowsPerPage);
    setCurrentPage(1); // Reset to first page
  }, []);

  // Navigate to specific row
  const goToRow = useCallback((rowNumber: number) => {
    const pageForRow = Math.ceil(rowNumber / rowsPerPage);
    setCurrentPage(pageForRow);
  }, [rowsPerPage]);

  // Handle drag start
  const handleDragStart = useCallback((e: React.DragEvent, globalIndex: number) => {
    setDraggedRowIndex(globalIndex);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/html", "");
  }, []);

  // Handle drag over
  const handleDragOver = useCallback((e: React.DragEvent, globalIndex: number) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    setDragOverIndex(globalIndex);
  }, []);

  // Handle drag leave
  const handleDragLeave = useCallback(() => {
    setDragOverIndex(null);
  }, []);

  // Handle drop
  const handleDrop = useCallback(async (e: React.DragEvent, dropIndex: number) => {
    e.preventDefault();
    setDragOverIndex(null);
    
    if (draggedRowIndex === null || draggedRowIndex === dropIndex || !token || !currentFile) {
      setDraggedRowIndex(null);
      return;
    }

    const draggedRow = allCsvRows[draggedRowIndex];
    const newPosition = calculateInsertPosition(dropIndex);

    try {
      const res = await fetch(`http://localhost:8080/files/${currentFile.id}/rows/${draggedRow.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify({
          position: newPosition,
          input_text: draggedRow.input_text,
        }),
      });

      const data = await res.json();
      if (data.status) {
        fetchCSVRows(currentFile, true); // Maintain current page
      } else {
        alert("Failed to reorder row: " + data.message);
      }
    } catch (err) {
      console.error("Error reordering row:", err);
      alert("Failed to reorder row");
    }
    
    setDraggedRowIndex(null);
  }, [draggedRowIndex, token, currentFile, allCsvRows, calculateInsertPosition, fetchCSVRows]);

  // Handle delete row
  const handleDeleteRow = useCallback(async (row: CSVRow, globalIndex: number) => {
    if (!token || !currentFile) return;

    if (!confirm("Are you sure you want to delete this row?")) return;

    try {
      const res = await fetch(`http://localhost:8080/files/${currentFile.id}/rows/${row.id}`, {
        method: "DELETE",
        headers: {
          Authorization: "Bearer " + token,
        },
      });

      if (!res.ok) {
        const text = await res.text();
        alert("Failed to delete row: " + text);
        return;
      }

      const data = await res.json();
      if (data.status) {
        fetchCSVRows(currentFile, true); // Maintain current page
      } else {
        alert("Failed to delete row: " + data.message);
      }
    } catch (err) {
      console.error("Error deleting row:", err);
      alert("Failed to delete row");
    }
  }, [token, currentFile, fetchCSVRows]);

  // Handle edit row with optimistic update
  const handleEditRow = useCallback(async (row: CSVRow, newText: string) => {
    if (!token || !currentFile || !newText.trim()) return;
    
    // Optimistic update
    setAllCsvRows(prev => prev.map(r => 
      r.id === row.id ? { ...r, input_text: newText.trim() } : r
    ));
    
    try {
      const res = await fetch(`http://localhost:8080/files/${currentFile.id}/rows/${row.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + token,
        },
        body: JSON.stringify({
          position: row.position,
          input_text: newText.trim(),
        }),
      });

      const data = await res.json();
      if (!data.status) {
        // Revert on failure
        alert("Failed to update row: " + data.message);
        fetchCSVRows(currentFile, true); // Maintain current page
      }
    } catch (err) {
      console.error("Error updating row:", err);
      alert("Failed to update row");
      fetchCSVRows(currentFile);
    }
  }, [token, currentFile, fetchCSVRows]);

  // If viewing CSV data
  if (currentFile) {
    return (
      <Box p={4}>
        <Box display="flex" alignItems="center" mb={2}>
          <IconButton onClick={handleBackToFiles} sx={{ mr: 2 }}>
            <ArrowBack />
          </IconButton>
          <Typography variant="h4">
            {currentFile.filename}
          </Typography>
          <Chip 
            label={`${totalRows} total rows`} 
            color="primary" 
            size="small" 
            sx={{ ml: 2 }} 
          />
          <Chip 
            label={`Page ${currentPage} of ${totalPages}`} 
            color="secondary" 
            size="small" 
            sx={{ ml: 1 }} 
          />
        </Box>

        {/* Pagination Controls */}
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Box display="flex" alignItems="center" gap={2}>
            <FormControl size="small" sx={{ minWidth: 120 }}>
              <InputLabel>Rows per page</InputLabel>
              <Select
                value={rowsPerPage}
                label="Rows per page"
                onChange={handleRowsPerPageChange}
              >
                <MenuItem value={50}>50</MenuItem>
                <MenuItem value={100}>100</MenuItem>
                <MenuItem value={200}>200</MenuItem>
                <MenuItem value={500}>500</MenuItem>
              </Select>
            </FormControl>
            
            <Typography variant="body2" color="textSecondary">
              Showing rows {startIndex + 1}-{Math.min(endIndex, totalRows)} of {totalRows}
            </Typography>
          </Box>
          
          <Pagination 
            count={totalPages}
            page={currentPage}
            onChange={handlePageChange}
            color="primary"
            size="large"
          />
        </Box>

        {loading ? (
          <Box display="flex" justifyContent="center" p={4}>
            <CircularProgress />
          </Box>
        ) : (
          <TableContainer component={Paper} sx={{ maxHeight: 600 }}>
            <Table stickyHeader>
              <TableHead>
                <TableRow>
                  <TableCell sx={{ width: 120 }}>Actions</TableCell>
                  <TableCell sx={{ width: 80 }}>Row #</TableCell>
                  <TableCell sx={{ width: 100 }}>Position</TableCell>
                  <TableCell>Content</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {currentPageRows.map((row, index) => {
                  const globalIndex = startIndex + index;
                  return (
                    <CSVRowComponent
                      key={row.id}
                      row={row}
                      index={index}
                      globalIndex={globalIndex}
                      draggedRowIndex={draggedRowIndex}
                      dragOverIndex={dragOverIndex}
                      onDragStart={handleDragStart}
                      onDragOver={handleDragOver}
                      onDragLeave={handleDragLeave}
                      onDrop={handleDrop}
                      onInsertRow={handleInsertRow}
                      onDeleteRow={handleDeleteRow}
                      onEditRow={handleEditRow}
                    />
                  );
                })}
                
                {/* Add row at end button - only show on last page */}
                {currentPage === totalPages && (
                  <TableRow>
                    <TableCell colSpan={4}>
                      <Button
                        startIcon={<Add />}
                        onClick={() => handleInsertRow(totalRows)}
                        variant="outlined"
                        fullWidth
                        sx={{ py: 2 }}
                      >
                        Add Row at End
                      </Button>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        )}

        {/* Bottom Pagination */}
        <Box display="flex" justifyContent="center" mt={2}>
          <Pagination 
            count={totalPages}
            page={currentPage}
            onChange={handlePageChange}
            color="primary"
            size="large"
            showFirstButton
            showLastButton
          />
        </Box>

        {/* Insert Row Dialog */}
        <Dialog 
          open={insertDialogOpen} 
          onClose={() => setInsertDialogOpen(false)}
          maxWidth="sm"
          fullWidth
        >
          <DialogTitle>
            Insert New Row at Position {insertPosition + 1}
          </DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              label="Row Content"
              fullWidth
              multiline
              rows={4}
              variant="outlined"
              value={newRowText}
              onChange={(e) => setNewRowText(e.target.value)}
              sx={{ mt: 2 }}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setInsertDialogOpen(false)}>
              Cancel
            </Button>
            <Button 
              onClick={confirmInsertRow} 
              variant="contained"
              disabled={!newRowText.trim()}
            >
              Insert Row
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    );
  }

  // File list view
  return (
    <Box p={4}>
      <Typography variant="h4" gutterBottom>
        CSV File Manager
      </Typography>

      <Paper sx={{ p: 2, mb: 3 }}>
        <Box display="flex" gap={2} alignItems="center">
          <Button variant="outlined" component="label">
            Select CSV
            <input
              type="file"
              accept=".csv"
              hidden
              onChange={handleFileChange}
            />
          </Button>

          <TextField
            label="Enter filename"
            value={filename}
            onChange={(e) => setFilename(e.target.value)}
          />

          <Button variant="contained" color="primary" onClick={handleUpload}>
            Upload
          </Button>
        </Box>
        {selectedFile && (
          <Typography variant="body2" sx={{ mt: 1, color: "text.secondary" }}>
            Selected: {selectedFile.name}
          </Typography>
        )}
      </Paper>

      <Typography variant="h6" gutterBottom>
        Uploaded Files
      </Typography>

      {loading ? (
        <Box display="flex" justifyContent="center" p={4}>
          <CircularProgress />
        </Box>
      ) : (
        <List>
          {files.map((file) => (
            <ListItem
              key={file.id}
              onClick={() => handleFileClick(file)}
              sx={{
                borderBottom: "1px solid #ddd",
                cursor: "pointer",
                "&:hover": { backgroundColor: "#f5f5f5" },
              }}
            >
              <ListItemText
                primary={file.filename}
                secondary={new Date(file.uploaded_at).toLocaleString()}
              />
            </ListItem>
          ))}
        </List>
      )}
    </Box>
  );
}