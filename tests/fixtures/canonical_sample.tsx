import React from "react";

// Canonical benchmark fixture for Charites static analysis parser
export function CanonicalDashboard() {
  return (
    <div className="min-h-screen bg-slate-50 p-8">
      <header className="mb-8 border-b pb-4">
        <h1 className="text-2xl font-bold text-slate-800">Analytics Center</h1>
        <p className="text-sm text-slate-500">Comprehensive metric monitoring and diagnostic view</p>
      </header>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-1</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-2</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-3</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-4</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-5</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-6</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-7</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-8</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-9</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-10</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-11</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-12</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-13</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-14</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-15</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-16</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-17</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-18</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-19</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-20</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-21</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-22</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-23</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-24</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-25</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-26</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-27</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-28</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-29</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-30</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-31</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-32</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-33</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-34</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-35</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-36</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-37</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-38</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-39</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-40</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-41</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-42</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-43</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-44</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-45</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-46</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-47</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-48</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-49</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-50</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-51</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-52</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-53</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-54</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-55</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-56</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-57</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-58</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-59</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-60</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-61</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-62</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-63</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-64</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-65</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-66</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-67</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-68</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-69</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-70</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-71</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-72</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-73</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-74</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-75</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-76</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-77</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-78</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-79</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-80</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-81</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-82</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-83</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-84</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-85</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-86</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-87</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-88</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-89</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-90</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-91</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-92</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-93</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-94</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-95</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-96</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-97</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-98</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-99</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-100</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-101</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-102</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-103</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-104</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-105</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-106</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-107</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-108</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-109</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-110</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-111</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-112</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-113</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-114</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-115</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-116</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-117</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-118</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-119</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>
        <div className={`p-4 rounded-xl border bg-white shadow-sm ${i % 2 === 0 ? "border-primary-light" : "border-slate-200"}`}>
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold text-slate-400">METRIC-120</span>
            <span className="inline-flex px-2 py-0.5 text-xs rounded bg-slate-100 text-slate-700">Active</span>
          </div>
          <p className="mt-2 text-xl font-bold text-slate-900">{`Value: ${i * 1234}`}</p>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1 text-xs rounded bg-primary text-white hover:bg-primary-hover" type="button">Inspect</button>
            <input type="checkbox" className="form-checkbox h-4 w-4 text-primary" defaultChecked />
          </div>
        </div>

      </div>
    </div>
  );
}
