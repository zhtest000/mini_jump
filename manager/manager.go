package manager

import (
	"net/http"
)

// Manager 管理页面处理器
type Manager struct{}

// NewManager 创建管理页面处理器
func NewManager() *Manager {
	return &Manager{}
}

// ServeManager 提供管理页面
func (m *Manager) ServeManager(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(managerHTML))
}

const managerHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MiniJump 规则管理</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            font-size: 28px;
            margin-bottom: 10px;
        }
        .content {
            padding: 30px;
        }
        .toolbar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            flex-wrap: wrap;
            gap: 10px;
        }
        .btn {
            padding: 10px 20px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.3s;
        }
        .btn-primary {
            background: #667eea;
            color: white;
        }
        .btn-primary:hover {
            background: #5568d3;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }
        .btn-success {
            background: #48bb78;
            color: white;
        }
        .btn-success:hover {
            background: #38a169;
        }
        .btn-danger {
            background: #f56565;
            color: white;
        }
        .btn-danger:hover {
            background: #e53e3e;
        }
        .btn-secondary {
            background: #718096;
            color: white;
        }
        .btn-secondary:hover {
            background: #4a5568;
        }
        .btn-small {
            padding: 5px 12px;
            font-size: 12px;
        }
        .table-container {
            overflow-x: auto;
            margin-top: 20px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            background: white;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #e2e8f0;
        }
        th {
            background: #f7fafc;
            font-weight: 600;
            color: #2d3748;
        }
        tr:hover {
            background: #f7fafc;
        }
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0,0,0,0.5);
            z-index: 1000;
            justify-content: center;
            align-items: center;
        }
        .modal.active {
            display: flex;
        }
        .modal-content {
            background: white;
            border-radius: 12px;
            padding: 30px;
            max-width: 600px;
            width: 90%;
            max-height: 90vh;
            overflow-y: auto;
        }
        .modal-header {
            font-size: 24px;
            font-weight: 600;
            margin-bottom: 20px;
            color: #2d3748;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #4a5568;
        }
        .form-group input,
        .form-group select,
        .form-group textarea {
            width: 100%;
            padding: 10px;
            border: 1px solid #cbd5e0;
            border-radius: 6px;
            font-size: 14px;
        }
        .form-group input:focus,
        .form-group select:focus,
        .form-group textarea:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        .form-actions {
            display: flex;
            justify-content: flex-end;
            gap: 10px;
            margin-top: 20px;
        }
        .alert {
            padding: 12px 16px;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        .alert-error {
            background: #fed7d7;
            color: #c53030;
            border: 1px solid #feb2b2;
        }
        .alert-success {
            background: #c6f6d5;
            color: #22543d;
            border: 1px solid #9ae6b4;
        }
        .alert-warning {
            background: #feebc8;
            color: #7c2d12;
            border: 1px solid #fbd38d;
        }
        .conflict-list {
            margin-top: 10px;
            padding: 10px;
            background: #fff5f5;
            border-radius: 6px;
        }
        .conflict-item {
            padding: 8px;
            margin: 5px 0;
            background: white;
            border-radius: 4px;
            border-left: 3px solid #f56565;
        }
        .badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 500;
        }
        .badge-301 { background: #4299e1; color: white; }
        .badge-302 { background: #48bb78; color: white; }
        .badge-307 { background: #ed8936; color: white; }
        .badge-4 { background: #9f7aea; color: white; }
        .status-expired {
            color: #a0aec0;
            text-decoration: line-through;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 MiniJump 跳转规则管理</h1>
            <p>轻量级 HTTP 跳转服务管理平台</p>
        </div>
        <div class="content">
            <div class="toolbar">
                <button class="btn btn-primary" onclick="openCreateModal()">+ 添加规则</button>
                <div>
                    <button class="btn btn-success" onclick="reloadConfig()">重新加载</button>
                    <button class="btn btn-secondary" onclick="saveConfig()">保存配置</button>
                </div>
            </div>
            <div id="alert-container"></div>
            <div class="table-container">
                <table id="rules-table">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>域名</th>
                            <th>路径</th>
                            <th>目标URL</th>
                            <th>类型</th>
                            <th>有效期</th>
                            <th>描述</th>
                            <th>操作</th>
                        </tr>
                    </thead>
                    <tbody id="rules-tbody">
                        <tr>
                            <td colspan="8" style="text-align: center; padding: 40px; color: #a0aec0;">
                                加载中...
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>

    <!-- 创建/编辑规则模态框 -->
    <div id="rule-modal" class="modal">
        <div class="modal-content">
            <div class="modal-header" id="modal-title">添加规则</div>
            <div id="modal-alert"></div>
            <form id="rule-form" onsubmit="saveRule(event)">
                <input type="hidden" id="rule-id">
                <div class="form-group">
                    <label>域名 *</label>
                    <input type="text" id="rule-domain" required placeholder="example.com">
                </div>
                <div class="form-group">
                    <label>路径（可选）</label>
                    <input type="text" id="rule-path" placeholder="/old/path">
                </div>
                <div class="form-group">
                    <label>目标URL *</label>
                    <input type="url" id="rule-target" required placeholder="https://example.com/new">
                </div>
                <div class="form-group">
                    <label>跳转类型 *</label>
                    <select id="rule-type" required>
                        <option value="301">301 - 永久重定向</option>
                        <option value="302">302 - 临时重定向</option>
                        <option value="307">307 - 临时重定向（保持方法）</option>
                        <option value="4">JavaScript 跳转</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>有效期（可选）</label>
                    <input type="datetime-local" id="rule-expires">
                </div>
                <div class="form-group">
                    <label>描述</label>
                    <textarea id="rule-description" rows="3" placeholder="规则描述"></textarea>
                </div>
                <div class="form-actions">
                    <button type="button" class="btn btn-secondary" onclick="closeModal()">取消</button>
                    <button type="submit" class="btn btn-primary">保存</button>
                </div>
            </form>
        </div>
    </div>

    <script>
        let currentEditingId = null;

        // 加载规则列表
        async function loadRules() {
            try {
                const response = await fetch('/api/rules');
                if (!response.ok) throw new Error('加载失败');
                const rules = await response.json();
                renderRules(rules);
            } catch (error) {
                showAlert('加载规则失败: ' + error.message, 'error');
            }
        }

        // 渲染规则列表
        function renderRules(rules) {
            const tbody = document.getElementById('rules-tbody');
            if (rules.length === 0) {
                tbody.innerHTML = '<tr><td colspan="8" style="text-align: center; padding: 40px; color: #a0aec0;">暂无规则</td></tr>';
                return;
            }
            tbody.innerHTML = rules.map(rule => {
                const expiresAt = rule.expires_at ? new Date(rule.expires_at).toLocaleString('zh-CN') : '永不过期';
                const isExpired = rule.expires_at && new Date(rule.expires_at) < new Date();
                const expiredClass = isExpired ? 'status-expired' : '';
                const typeNames = {301: '301', 302: '302', 307: '307', 4: 'JS'};
                const path = rule.path || '<span style="color: #a0aec0;">（域名级别）</span>';
                const typeName = typeNames[rule.type] || rule.type;
                const description = rule.description || '-';
                return '<tr class="' + expiredClass + '">' +
                    '<td>' + rule.id + '</td>' +
                    '<td>' + rule.domain + '</td>' +
                    '<td>' + path + '</td>' +
                    '<td>' + rule.target + '</td>' +
                    '<td><span class="badge badge-' + rule.type + '">' + typeName + '</span></td>' +
                    '<td>' + expiresAt + '</td>' +
                    '<td>' + description + '</td>' +
                    '<td>' +
                    '<button class="btn btn-primary btn-small" onclick="editRule(\'' + rule.id + '\')">编辑</button> ' +
                    '<button class="btn btn-danger btn-small" onclick="deleteRule(\'' + rule.id + '\')">删除</button>' +
                    '</td>' +
                    '</tr>';
            }).join('');
        }

        // 打开创建模态框
        function openCreateModal() {
            currentEditingId = null;
            document.getElementById('modal-title').textContent = '添加规则';
            document.getElementById('rule-form').reset();
            document.getElementById('rule-id').value = '';
            document.getElementById('modal-alert').innerHTML = '';
            document.getElementById('rule-modal').classList.add('active');
        }

        // 编辑规则
        async function editRule(id) {
            try {
                const response = await fetch('/api/rules/' + id);
                if (!response.ok) throw new Error('加载失败');
                const rule = await response.json();
                
                currentEditingId = id;
                document.getElementById('modal-title').textContent = '编辑规则';
                document.getElementById('rule-id').value = rule.id;
                document.getElementById('rule-domain').value = rule.domain;
                document.getElementById('rule-path').value = rule.path || '';
                document.getElementById('rule-target').value = rule.target;
                document.getElementById('rule-type').value = rule.type;
                document.getElementById('rule-description').value = rule.description || '';
                
                if (rule.expires_at) {
                    const date = new Date(rule.expires_at);
                    const localDate = new Date(date.getTime() - date.getTimezoneOffset() * 60000);
                    document.getElementById('rule-expires').value = localDate.toISOString().slice(0, 16);
                } else {
                    document.getElementById('rule-expires').value = '';
                }
                
                document.getElementById('modal-alert').innerHTML = '';
                document.getElementById('rule-modal').classList.add('active');
            } catch (error) {
                showAlert('加载规则失败: ' + error.message, 'error');
            }
        }

        // 保存规则
        async function saveRule(event) {
            event.preventDefault();
            const alertDiv = document.getElementById('modal-alert');
            alertDiv.innerHTML = '';

            const ruleData = {
                domain: document.getElementById('rule-domain').value.trim(),
                path: document.getElementById('rule-path').value.trim(),
                target: document.getElementById('rule-target').value.trim(),
                type: parseInt(document.getElementById('rule-type').value),
                description: document.getElementById('rule-description').value.trim()
            };

            const expiresValue = document.getElementById('rule-expires').value;
            if (expiresValue) {
                const expiresDate = new Date(expiresValue);
                ruleData.expires_at = expiresDate.toISOString();
            }

            try {
                let response;
                if (currentEditingId) {
                    ruleData.id = currentEditingId;
                    response = await fetch('/api/rules/' + currentEditingId, {
                        method: 'PUT',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify(ruleData)
                    });
                } else {
                    response = await fetch('/api/rules', {
                        method: 'POST',
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify(ruleData)
                    });
                }

                const result = await response.json();
                
                if (!response.ok) {
                    if (response.status === 409) {
                        // 冲突错误
                        let conflictHTML = '<div class="alert alert-error"><strong>规则冲突！</strong><br>' + result.error + '</div>';
                        if (result.conflicts && result.conflicts.length > 0) {
                            conflictHTML += '<div class="conflict-list"><strong>冲突的规则：</strong>';
                            result.conflicts.forEach(conflict => {
                                conflictHTML += '<div class="conflict-item">';
                                const conflictPath = conflict.path || '(域名级别)';
                                conflictHTML += 'ID: ' + conflict.id + ', 域名: ' + conflict.domain + ', 路径: ' + conflictPath;
                                conflictHTML += '</div>';
                            });
                            conflictHTML += '</div>';
                        }
                        alertDiv.innerHTML = conflictHTML;
                    } else {
                        alertDiv.innerHTML = '<div class="alert alert-error">保存失败: ' + (result.error || '未知错误') + '</div>';
                    }
                    return;
                }

                closeModal();
                showAlert('规则保存成功', 'success');
                loadRules();
            } catch (error) {
                alertDiv.innerHTML = '<div class="alert alert-error">保存失败: ' + error.message + '</div>';
            }
        }

        // 删除规则
        async function deleteRule(id) {
            if (!confirm('确定要删除这条规则吗？')) return;
            
            try {
                const response = await fetch('/api/rules/' + id, {method: 'DELETE'});
                if (!response.ok) throw new Error('删除失败');
                showAlert('规则删除成功', 'success');
                loadRules();
            } catch (error) {
                showAlert('删除失败: ' + error.message, 'error');
            }
        }

        // 关闭模态框
        function closeModal() {
            document.getElementById('rule-modal').classList.remove('active');
        }

        // 重新加载配置
        async function reloadConfig() {
            try {
                const response = await fetch('/api/reload', {method: 'POST'});
                if (!response.ok) throw new Error('重新加载失败');
                showAlert('配置重新加载成功', 'success');
                loadRules();
            } catch (error) {
                showAlert('重新加载失败: ' + error.message, 'error');
            }
        }

        // 保存配置
        async function saveConfig() {
            try {
                const response = await fetch('/api/save', {method: 'POST'});
                if (!response.ok) throw new Error('保存失败');
                showAlert('配置保存成功', 'success');
            } catch (error) {
                showAlert('保存失败: ' + error.message, 'error');
            }
        }

        // 显示提示信息
        function showAlert(message, type) {
            const container = document.getElementById('alert-container');
            const alert = document.createElement('div');
            alert.className = 'alert alert-' + type;
            alert.textContent = message;
            container.innerHTML = '';
            container.appendChild(alert);
            setTimeout(() => {
                alert.remove();
            }, 3000);
        }

        // 点击模态框外部关闭
        document.getElementById('rule-modal').addEventListener('click', function(e) {
            if (e.target === this) {
                closeModal();
            }
        });

        // 页面加载时获取规则列表
        loadRules();
    </script>
</body>
</html>
`
