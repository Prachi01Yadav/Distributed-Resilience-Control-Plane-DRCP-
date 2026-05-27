// ==========================================
// DRCP - Distributed Resilience Control Plane
// Frontend Application Logic
// ==========================================

const API = '/api/v1';

// Cached data
let allServices = [];
let allContracts = [];
let allIncidents = [];

// ==========================================
// Initialization
// ==========================================
document.addEventListener('DOMContentLoaded', () => {
    lucide.createIcons();
    setupNavigation();
    loadAllData();
});

async function loadAllData() {
    await Promise.all([
        fetchServices(),
        fetchContracts(),
        fetchIncidents()
    ]);
    updateDashboardStats();
}

// ==========================================
// Navigation
// ==========================================
function setupNavigation() {
    document.querySelectorAll('.nav-item').forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const viewName = e.currentTarget.getAttribute('data-view');
            switchView(viewName);
        });
    });
}

function switchView(viewName) {
    // Update nav
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
    const activeNav = document.querySelector(`[data-view="${viewName}"]`);
    if (activeNav) activeNav.classList.add('active');

    // Update views
    document.querySelectorAll('.view-section').forEach(v => {
        v.classList.remove('active');
        v.style.display = 'none';
    });
    const target = document.getElementById(`view-${viewName}`);
    if (target) {
        target.style.display = 'block';
        // Re-trigger animation
        target.classList.remove('active');
        void target.offsetWidth;
        target.classList.add('active');
    }

    // Refresh data for the view
    if (viewName === 'services') renderServicesView();
    if (viewName === 'contracts') renderContractsView();
    if (viewName === 'incidents') renderIncidentsView();
    if (viewName === 'dashboard') {
        updateDashboardStats();
        renderDashboardTable();
    }

    lucide.createIcons();
}

// ==========================================
// Data Fetching
// ==========================================
async function fetchServices() {
    try {
        const res = await fetch(`${API}/services`);
        if (!res.ok) throw new Error('API error');
        allServices = await res.json() || [];
    } catch (e) {
        console.error('Failed to fetch services:', e);
        allServices = [];
    }
    renderDashboardTable();
}

async function fetchContracts() {
    try {
        const res = await fetch(`${API}/contracts`);
        if (!res.ok) throw new Error('API error');
        allContracts = await res.json() || [];
    } catch (e) {
        console.error('Failed to fetch contracts:', e);
        allContracts = [];
    }
}

async function fetchIncidents() {
    try {
        const res = await fetch(`${API}/incidents`);
        if (!res.ok) throw new Error('API error');
        allIncidents = await res.json() || [];
    } catch (e) {
        console.error('Failed to fetch incidents:', e);
        allIncidents = [];
    }
}

// ==========================================
// Dashboard
// ==========================================
function updateDashboardStats() {
    document.getElementById('stat-total-services').textContent = allServices.length;
    document.getElementById('stat-active-slas').textContent = allContracts.filter(c => c.is_active).length;
    const openIncidents = allIncidents.filter(i => i.status === 'OPEN').length;
    document.getElementById('stat-open-incidents').textContent = openIncidents;
}

function renderDashboardTable() {
    const tbody = document.getElementById('dashboard-services-tbody');
    if (!allServices.length) {
        tbody.innerHTML = `<tr><td colspan="6" class="empty-cell">No services registered. Click "New Service" to begin.</td></tr>`;
        return;
    }
    tbody.innerHTML = allServices.slice(0, 5).map(s => serviceRow(s, false)).join('');
    lucide.createIcons();
}

// ==========================================
// Services View
// ==========================================
function renderServicesView() {
    const tbody = document.getElementById('services-tbody');
    if (!allServices.length) {
        tbody.innerHTML = `<tr><td colspan="7" class="empty-cell">No services registered. Click "Register Service" to begin.</td></tr>`;
        return;
    }
    tbody.innerHTML = allServices.map(s => {
        const initial = s.owner ? s.owner.charAt(0).toUpperCase() : '?';
        const contractCount = allContracts.filter(c => c.service_id === s.id).length;
        return `<tr>
            <td>
                <div style="font-weight:600;">${esc(s.name)}</div>
                <div style="font-size:0.75rem;color:var(--text-muted);">${s.id.substring(0,8)}</div>
            </td>
            <td>
                <div style="display:flex;align-items:center;gap:0.75rem;">
                    <div class="owner-avatar">${initial}</div>
                    <span style="font-weight:500;">${esc(s.owner)}</span>
                </div>
            </td>
            <td style="color:var(--text-secondary);max-width:200px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;">${esc(s.description || 'No description')}</td>
            <td><span class="badge badge-success"><i data-lucide="shield-check" style="width:14px;height:14px;"></i> Protected</span></td>
            <td>
                <div style="display:flex;align-items:center;gap:0.5rem;">
                    <div class="pulse-dot"></div> <span style="font-weight:600;">100%</span>
                </div>
            </td>
            <td>${formatDate(s.created_at)}</td>
            <td>
                <div style="display:flex;gap:0.5rem;">
                    <button class="btn btn-secondary btn-sm" onclick="openContractForService('${s.id}','${esc(s.name)}')">+ SLA</button>
                    <button class="btn btn-secondary btn-sm" onclick="showToast('Details view coming in Phase 2')" title="View details"><i data-lucide="eye" style="width:14px;height:14px;"></i></button>
                </div>
            </td>
        </tr>`;
    }).join('');
    lucide.createIcons();
}

function serviceRow(s, showDesc) {
    const initial = s.owner ? s.owner.charAt(0).toUpperCase() : '?';
    return `<tr>
        <td>
            <div style="font-weight:600;">${esc(s.name)}</div>
            <div style="font-size:0.75rem;color:var(--text-muted);">${s.id.substring(0,8)}</div>
        </td>
        <td>
            <div style="display:flex;align-items:center;gap:0.75rem;">
                <div class="owner-avatar">${initial}</div>
                <span style="font-weight:500;">${esc(s.owner)}</span>
            </div>
        </td>
        <td><span class="badge badge-success"><i data-lucide="shield-check" style="width:14px;height:14px;"></i> Protected</span></td>
        <td>
            <div style="display:flex;align-items:center;gap:0.5rem;">
                <div class="pulse-dot"></div> <span style="font-weight:600;">100%</span>
            </div>
        </td>
        <td>${formatDate(s.created_at)}</td>
        <td>
            <button class="icon-btn" onclick="showToast('Service options coming in Phase 2')"><i data-lucide="more-horizontal"></i></button>
        </td>
    </tr>`;
}

// ==========================================
// SLA Contracts View
// ==========================================
function renderContractsView() {
    const tbody = document.getElementById('contracts-tbody');
    if (!allContracts.length) {
        tbody.innerHTML = `<tr><td colspan="5" class="empty-cell">No SLA contracts created yet. Click "New Contract" to create one.</td></tr>`;
        return;
    }
    tbody.innerHTML = allContracts.map(c => {
        const svc = allServices.find(s => s.id === c.service_id);
        const svcName = svc ? svc.name : c.service_id.substring(0, 8);
        const statusBadge = c.is_active
            ? '<span class="badge badge-success">Active</span>'
            : '<span class="badge badge-muted">Inactive</span>';
        return `<tr>
            <td>
                <div style="font-weight:600;">${c.id.substring(0,8)}</div>
            </td>
            <td>
                <div style="font-weight:500;">${esc(svcName)}</div>
                <div style="font-size:0.75rem;color:var(--text-muted);">${c.service_id.substring(0,8)}</div>
            </td>
            <td><div class="policy-code">${esc(c.policy)}</div></td>
            <td>${statusBadge}</td>
            <td>${formatDate(c.created_at)}</td>
        </tr>`;
    }).join('');
    lucide.createIcons();
}

function openContractForService(serviceId, serviceName) {
    openModal('contract-modal');
    populateServiceDropdown();
    document.getElementById('contract-service-id').value = serviceId;
}

function populateServiceDropdown() {
    const select = document.getElementById('contract-service-id');
    select.innerHTML = '<option value="">Select a service...</option>';
    allServices.forEach(s => {
        select.innerHTML += `<option value="${s.id}">${esc(s.name)}</option>`;
    });
}

// ==========================================
// Incidents View
// ==========================================
function renderIncidentsView() {
    // Stats
    document.getElementById('incidents-total').textContent = allIncidents.length;
    document.getElementById('incidents-open').textContent = allIncidents.filter(i => i.status === 'OPEN').length;
    document.getElementById('incidents-anchored').textContent = allIncidents.filter(i => i.status === 'ANCHORED').length;

    const tbody = document.getElementById('incidents-tbody');
    if (!allIncidents.length) {
        tbody.innerHTML = `<tr><td colspan="7" class="empty-cell">No incidents recorded. Your services are running within SLA bounds.</td></tr>`;
        return;
    }
    tbody.innerHTML = allIncidents.map(inc => {
        let statusBadge;
        switch (inc.status) {
            case 'OPEN': statusBadge = '<span class="badge badge-danger">Open</span>'; break;
            case 'RESOLVED': statusBadge = '<span class="badge badge-success">Resolved</span>'; break;
            case 'ANCHORED': statusBadge = '<span class="badge badge-info">Anchored</span>'; break;
            default: statusBadge = '<span class="badge badge-muted">' + esc(inc.status) + '</span>';
        }
        const txHash = inc.blockchain_tx_hash
            ? `<span class="policy-code" title="${esc(inc.blockchain_tx_hash)}">${inc.blockchain_tx_hash.substring(0,12)}...</span>`
            : '<span style="color:var(--text-muted);">Pending</span>';
        return `<tr>
            <td><div style="font-weight:600;">${inc.id.substring(0,8)}</div></td>
            <td>${inc.service_id.substring(0,8)}</td>
            <td><span style="font-weight:600;color:${inc.error_rate > 0.01 ? 'var(--danger)' : 'var(--success)'};">${(inc.error_rate * 100).toFixed(2)}%</span></td>
            <td><span style="font-weight:600;color:${inc.p99_latency > 500 ? 'var(--danger)' : 'var(--success)'};">${inc.p99_latency.toFixed(0)}ms</span></td>
            <td>${statusBadge}</td>
            <td>${txHash}</td>
            <td>${formatDate(inc.created_at)}</td>
        </tr>`;
    }).join('');
    lucide.createIcons();
}

// ==========================================
// Modals
// ==========================================
function openModal(id) {
    document.getElementById(id).classList.add('active');
    if (id === 'contract-modal') populateServiceDropdown();
    lucide.createIcons();
}

function closeModal(id) {
    document.getElementById(id).classList.remove('active');
}

// ==========================================
// Form Submissions
// ==========================================
async function submitService(event) {
    event.preventDefault();
    const name = document.getElementById('svc-name').value;
    const owner = document.getElementById('svc-owner').value;
    const desc = document.getElementById('svc-desc').value;

    try {
        const res = await fetch(`${API}/services`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, owner, description: desc })
        });
        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || 'Failed to create service');
        }
        closeModal('service-modal');
        document.getElementById('register-form').reset();
        showToast(`Service "${name}" registered successfully.`);
        await fetchServices();
        updateDashboardStats();
    } catch (e) {
        showToast(e.message, 'error');
    }
}

async function submitContract(event) {
    event.preventDefault();
    const serviceId = document.getElementById('contract-service-id').value;
    const policy = document.getElementById('contract-policy').value;

    if (!serviceId) {
        showToast('Please select a service.', 'error');
        return;
    }

    try {
        const res = await fetch(`${API}/services/${serviceId}/contracts`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ policy })
        });
        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || 'Failed to create contract');
        }
        closeModal('contract-modal');
        document.getElementById('contract-form').reset();
        showToast('SLA Contract created successfully.');
        await fetchContracts();
        updateDashboardStats();
    } catch (e) {
        showToast(e.message, 'error');
    }
}

// ==========================================
// Toast System
// ==========================================
function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = 'toast';
    if (type === 'error') toast.style.borderLeftColor = 'var(--danger)';
    const iconName = type === 'error' ? 'alert-circle' : 'check-circle';
    const iconColor = type === 'error' ? 'var(--danger)' : 'var(--accent-primary)';
    toast.innerHTML = `<i data-lucide="${iconName}" style="width:18px;height:18px;color:${iconColor};flex-shrink:0;"></i> ${esc(message)}`;
    container.appendChild(toast);
    lucide.createIcons();
    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateY(10px)';
        toast.style.transition = 'all 0.3s ease';
        setTimeout(() => toast.remove(), 300);
    }, 3500);
}

function toggleNotifications() {
    const openCount = allIncidents.filter(i => i.status === 'OPEN').length;
    if (openCount > 0) {
        showToast(`You have ${openCount} open incident(s) requiring attention.`, 'error');
    } else {
        showToast('All clear. No open incidents.');
    }
}

// ==========================================
// Helpers
// ==========================================
function formatDate(dateStr) {
    if (!dateStr) return '-';
    return new Date(dateStr).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });
}

function esc(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}
