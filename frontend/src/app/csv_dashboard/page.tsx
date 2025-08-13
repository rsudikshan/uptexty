"use client";
import { useState, useEffect, ChangeEvent } from "react";
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
} from "@mui/material";
import {
  ArrowBack,
  Add,
  Delete,
  DragIndicator,
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

export default function FilesPage() {
  const [files, setFiles] = useState<FileItem[]>([]);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [filename, setFilename] = useState("");
  const [loading, setLoading] = useState(false);
  
  // New state for CSV data display
  const [currentFile, setCurrentFile] = useState<FileItem | null>(null);
  const [csvRows, setCsvRows] = useState<CSVRow[]>([]);
  const [insertDialogOpen, setInsertDialogOpen] = useState(false);
  const [insertPosition, setInsertPosition] = useState(0);
  const [newRowText, setNewRowText] = useState("");
  const [draggedRowIndex, setDraggedRowIndex] = useState<number | null>(null);
  const [dragOverIndex, setDragOverIndex] = useState<number | null>(null);

  const token = localStorage.getItem("jwt");

  // Fetch uploaded files
  const fetchFiles = async () => {
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
  };

  useEffect(() => {
    fetchFiles();
  }, []);

  // Handle file input
  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  };

  // Upload CSV file
  const handleUpload = async () => {
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
  };

  // Fetch CSV rows for display
  const fetchCSVRows = async (file: FileItem) => {
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
        // Sort rows by position to maintain order
        const sortedRows = (data.body || []).sort((a: CSVRow, b: CSVRow) => a.position - b.position);
        setCsvRows(sortedRows);
        setCurrentFile(file);
      } else {
        alert("Failed to load CSV data: " + data.message);
      }
    } catch (err) {
      console.error("Error fetching CSV data:", err);
      alert("Failed to load CSV data");
    } finally {
      setLoading(false);
    }
  };

  // Handle file click to display CSV data
  const handleFileClick = (file: FileItem) => {
    fetchCSVRows(file);
  };

  // Handle back to file list
  const handleBackToFiles = () => {
    setCurrentFile(null);
    setCsvRows([]);
  };

  // Calculate position for new row insertion
  const calculateInsertPosition = (index: number): number => {
    if (csvRows.length === 0) return 1.0;
    
    if (index === 0) {
      // Insert at beginning
      return csvRows[0].position / 2;
    } else if (index >= csvRows.length) {
      // Insert at end
      return csvRows[csvRows.length - 1].position + 1.0;
    } else {
      // Insert between two rows
      const prevPosition = csvRows[index - 1].position;
      const nextPosition = csvRows[index].position;
      return (prevPosition + nextPosition) / 2;
    }
  };

  // Handle insert row dialog
  const handleInsertRow = (index: number) => {
    setInsertPosition(index);
    setNewRowText("");
    setInsertDialogOpen(true);
  };

  // Confirm insert row
  const confirmInsertRow = async () => {
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
        // Refresh CSV data
        fetchCSVRows(currentFile);
        setInsertDialogOpen(false);
      } else {
        alert("Failed to insert row: " + data.message);
      }
    } catch (err) {
      console.error("Error inserting row:", err);
      alert("Failed to insert row");
    }
  };

  // Handle drag start
  const handleDragStart = (e: React.DragEvent, index: number) => {
    setDraggedRowIndex(index);
    e.dataTransfer.effectAllowed = "move";
    e.dataTransfer.setData("text/html", "");
  };

  // Handle drag over
  const handleDragOver = (e: React.DragEvent, index: number) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
    setDragOverIndex(index);
  };

  // Handle drag leave
  const handleDragLeave = () => {
    setDragOverIndex(null);
  };

  // Handle drop
  const handleDrop = async (e: React.DragEvent, dropIndex: number) => {
    e.preventDefault();
    setDragOverIndex(null);
    
    if (draggedRowIndex === null || draggedRowIndex === dropIndex || !token || !currentFile) {
      setDraggedRowIndex(null);
      return;
    }

    const draggedRow = csvRows[draggedRowIndex];
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
        // Refresh CSV data
        fetchCSVRows(currentFile);
      } else {
        alert("Failed to reorder row: " + data.message);
      }
    } catch (err) {
      console.error("Error reordering row:", err);
      alert("Failed to reorder row");
    }
    
    setDraggedRowIndex(null);
  };

  // Handle delete row
    const handleDeleteRow = async (row: CSVRow, index: number) => {
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
        const text = await res.text(); // read plain text if error
        alert("Failed to delete row: " + text);
        return;
      }

      const data = await res.json(); // now safe, should be JSON for success
      if (data.status) {
        fetchCSVRows(currentFile);
      } else {
        alert("Failed to delete row: " + data.message);
      }
    } catch (err) {
      console.error("Error deleting row:", err);
      alert("Failed to delete row");
    }
  };

  // Handle edit row
  const handleEditRow = async (row: CSVRow, newText: string) => {
    if (!token || !currentFile || !newText.trim()) return;
    
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
      if (data.status) {
        // Refresh CSV data
        fetchCSVRows(currentFile);
      } else {
        alert("Failed to update row: " + data.message);
      }
    } catch (err) {
      console.error("Error updating row:", err);
      alert("Failed to update row");
    }
  };

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
            label={`${csvRows.length} rows`} 
            color="primary" 
            size="small" 
            sx={{ ml: 2 }} 
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
                {csvRows.map((row, index) => (
                  <TableRow
                    key={row.id}
                    draggable
                    onDragStart={(e) => handleDragStart(e, index)}
                    onDragOver={(e) => handleDragOver(e, index)}
                    onDragLeave={handleDragLeave}
                    onDrop={(e) => handleDrop(e, index)}
                    sx={{
                      cursor: "move",
                      "&:hover": { backgroundColor: "#f5f5f5" },
                      backgroundColor: 
                        draggedRowIndex === index ? "#e3f2fd" : 
                        dragOverIndex === index ? "#f3e5f5" : "inherit",
                      borderLeft: dragOverIndex === index ? "3px solid #9c27b0" : "none",
                    }}
                  >
                    <TableCell>
                      <Box display="flex" gap={1}>
                        <IconButton size="small" sx={{ cursor: "grab" }}>
                          <DragIndicator />
                        </IconButton>
                        <IconButton 
                          size="small" 
                          color="primary"
                          onClick={() => handleInsertRow(index)}
                          title="Insert row before this one"
                        >
                          <Add />
                        </IconButton>
                        <IconButton 
                          size="small" 
                          color="error"
                          onClick={() => handleDeleteRow(row, index)}
                          title="Delete this row"
                        >
                          <Delete />
                        </IconButton>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="textSecondary">
                        {index + 1}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Typography variant="body2" color="textSecondary">
                        {row.position.toFixed(3)}
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <TextField
                        fullWidth
                        multiline
                        variant="outlined"
                        value={row.input_text}
                        onChange={(e) => {
                          // Update local state immediately for responsive UI
                          setCsvRows(prev => prev.map(r => 
                            r.id === row.id ? { ...r, input_text: e.target.value } : r
                          ));
                        }}
                        onBlur={(e) => {
                          // Save to backend on blur
                          if (e.target.value !== row.input_text) {
                            handleEditRow(row, e.target.value);
                          }
                        }}
                        size="small"
                        sx={{ minWidth: 300 }}
                      />
                    </TableCell>
                  </TableRow>
                ))}
                {/* Add row at end button */}
                <TableRow>
                  <TableCell colSpan={4}>
                    <Button
                      startIcon={<Add />}
                      onClick={() => handleInsertRow(csvRows.length)}
                      variant="outlined"
                      fullWidth
                      sx={{ py: 2 }}
                    >
                      Add Row at End
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        )}

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

  // File list view (original)
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
          {files.map((file, index) => (
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