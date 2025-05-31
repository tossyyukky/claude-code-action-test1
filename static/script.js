let currentTaskId = null;

function showTaskDetail(taskId) {
    fetch(`/api/tasks/${taskId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('タスクの取得に失敗しました');
            }
            return response.json();
        })
        .then(task => {
            document.getElementById('modalTaskTitle').textContent = task.title;
            document.getElementById('modalTaskDescription').textContent = task.description || '詳細情報なし';
            
            const deadline = new Date(task.deadline);
            document.getElementById('modalTaskDeadline').textContent = 
                deadline.toLocaleDateString('ja-JP', {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit'
                });
            
            document.getElementById('modalTaskStatus').textContent = task.completed ? '完了' : '未完了';
            
            const created = new Date(task.created_at);
            document.getElementById('modalTaskCreated').textContent = 
                created.toLocaleDateString('ja-JP', {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit'
                });
            
            currentTaskId = task.id;
            
            const completeBtn = document.getElementById('completeTaskBtn');
            const uncompleteBtn = document.getElementById('uncompleteTaskBtn');
            
            if (task.completed) {
                completeBtn.style.display = 'none';
                uncompleteBtn.style.display = 'inline-block';
            } else {
                completeBtn.style.display = 'inline-block';
                uncompleteBtn.style.display = 'none';
            }
            
            const modal = new bootstrap.Modal(document.getElementById('taskModal'));
            modal.show();
        })
        .catch(error => {
            console.error('Error:', error);
            alert('タスクの詳細取得に失敗しました: ' + error.message);
        });
}

function toggleTaskStatus(taskId, completed) {
    const button = event.target;
    const originalText = button.textContent;
    
    button.disabled = true;
    button.textContent = '処理中...';
    
    fetch(`/api/tasks/${taskId}/status`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            completed: completed
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('ステータスの更新に失敗しました');
        }
        return response.json();
    })
    .then(data => {
        if (data.status === 'success') {
            location.reload();
        } else {
            throw new Error('サーバーエラーが発生しました');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('ステータスの更新に失敗しました: ' + error.message);
        button.disabled = false;
        button.textContent = originalText;
    });
}

function toggleTaskStatusFromModal(completed) {
    if (!currentTaskId) return;
    
    const completeBtn = document.getElementById('completeTaskBtn');
    const uncompleteBtn = document.getElementById('uncompleteTaskBtn');
    
    completeBtn.disabled = true;
    uncompleteBtn.disabled = true;
    
    fetch(`/api/tasks/${currentTaskId}/status`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            completed: completed
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('ステータスの更新に失敗しました');
        }
        return response.json();
    })
    .then(data => {
        if (data.status === 'success') {
            bootstrap.Modal.getInstance(document.getElementById('taskModal')).hide();
            location.reload();
        } else {
            throw new Error('サーバーエラーが発生しました');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('ステータスの更新に失敗しました: ' + error.message);
        completeBtn.disabled = false;
        uncompleteBtn.disabled = false;
    });
}

document.addEventListener('DOMContentLoaded', function() {
    const completeBtn = document.getElementById('completeTaskBtn');
    const uncompleteBtn = document.getElementById('uncompleteTaskBtn');
    
    completeBtn.addEventListener('click', function() {
        toggleTaskStatusFromModal(true);
    });
    
    uncompleteBtn.addEventListener('click', function() {
        toggleTaskStatusFromModal(false);
    });
    
    document.getElementById('taskModal').addEventListener('hidden.bs.modal', function() {
        currentTaskId = null;
    });
});