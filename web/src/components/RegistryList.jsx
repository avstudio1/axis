const RegistryList = ({ registry, selectedIndex, mode, registryRef, getTagStyles }) => (
    <div ref={registryRef} className="flex-1 space-y-1 overflow-y-auto scrollbar-hide p-2 pb-2">
        {registry.map((item, i) => {
            const tagLabel = (item.type === 'keep')
                ? (item.status || 'Pending')
                : item.type;
            return (
            <div key={item.id} className={`p-2 border transition-all ${i === selectedIndex && mode === 'MANUAL' ? 'bg-emerald-950/30 border-emerald-500 text-emerald-300' : 'border-transparent text-gray-600'}`}>
                <div className="flex justify-between text-xs font-bold">
                    <span>{item.title}</span>
                    <span className={`text-[9px] uppercase px-2 py-0.5 rounded-full border ${getTagStyles(tagLabel)}`}>{tagLabel}</span>
                </div>
                <div className="text-[10px] truncate italic">{item.snippet || 'No content preview.'}</div>
            </div>
            );
        })}
        <div className="h-2"></div>
    </div>
);

export default RegistryList;
