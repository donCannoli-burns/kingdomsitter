// Reset only the two user-facing Kingdomsitter preferences.
void main() {
    set_property("_sidecar", "true");
    set_property("sidecar", "false");
    print("Kingdomsitter defaults restored: _sidecar=true, sidecar=false", "navy");
}
