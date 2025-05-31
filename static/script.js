let currentTaskId = null;

function toggleTask(taskId) {
    fetch(`/task/${taskId}/toggle`, {
        method: 'POST'
    })
    .then(response => {
        if (response.ok) {
            location.reload();
        } else {
            alert('エラーが発生しました');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('エラーが発生しました');
    });
}

function showTaskDetails(taskId) {
    currentTaskId = taskId;
    
    fetch(`/task/${taskId}`)
    .then(response => response.json())
    .then(task => {
        document.getElementById('modalTitle').textContent = task.title;
        document.getElementById('modalDescription').textContent = task.description || '詳細情報なし';
        
        const deadline = new Date(task.deadline);
        document.getElementById('modalDeadline').textContent = deadline.toLocaleString('ja-JP');
        document.getElementById('modalStatus').textContent = task.completed ? '完了' : '未完了';
        
        const toggleBtn = document.getElementById('modalToggleBtn');
        toggleBtn.textContent = task.completed ? '未完了にする' : '完了にする';
        toggleBtn.className = task.completed ? 'btn btn-secondary' : 'btn btn-success';
        
        const modal = new bootstrap.Modal(document.getElementById('taskModal'));
        modal.show();
    })
    .catch(error => {
        console.error('Error:', error);
        alert('タスクの詳細を取得できませんでした');
    });
}

function toggleTaskFromModal() {
    if (currentTaskId) {
        fetch(`/task/${currentTaskId}/toggle`, {
            method: 'POST'
        })
        .then(response => {
            if (response.ok) {
                const modal = bootstrap.Modal.getInstance(document.getElementById('taskModal'));
                modal.hide();
                location.reload();
            } else {
                alert('エラーが発生しました');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('エラーが発生しました');
        });
    }
}

// Close modal when clicking outside
document.addEventListener('DOMContentLoaded', function() {
    const taskModal = document.getElementById('taskModal');
    taskModal.addEventListener('hidden.bs.modal', function () {
        currentTaskId = null;
    });
});