import { useEffect, useRef } from "react";
import type { Awareness } from "y-protocols/awareness";
import type * as Y from "yjs";
import { EditorState } from "@codemirror/state";
import { EditorView, basicSetup } from "codemirror";
import { yCollab } from "y-codemirror.next";
import type { RemoteUser } from "../hooks/useCollaborativeEditor";
import "./CollaborativeMarkdownEditor.css";

export type CollaborativeMarkdownEditorProps = {
  yDoc: Y.Doc | null;
  awareness: Awareness | null;
  isConnected: boolean;
  remoteUsers: RemoteUser[];
  className?: string;
};

export const CollaborativeMarkdownEditor = ({
  yDoc,
  awareness,
  isConnected,
  remoteUsers,
  className = "",
}: CollaborativeMarkdownEditorProps) => {
  const editorContainerRef = useRef<HTMLDivElement | null>(null);
  const viewRef = useRef<EditorView | null>(null);

  useEffect(() => {
    if (!editorContainerRef.current || !yDoc || !awareness) {
      return;
    }

    const yText = yDoc.getText("content");
    const state = EditorState.create({
      doc: yText.toString(),
      extensions: [basicSetup, EditorView.lineWrapping, yCollab(yText, awareness)],
    });

    const view = new EditorView({
      state,
      parent: editorContainerRef.current,
    });
    viewRef.current = view;

    return () => {
      view.destroy();
      viewRef.current = null;
    };
  }, [awareness, yDoc]);

  return (
    <div className={`collaborative-editor ${className}`}>
      <div className="editor-header">
        <div className="status">
          <span
            className={`connection-status ${
              isConnected ? "connected" : "disconnected"
            }`}
          >
            {isConnected ? "Connected" : "Disconnected"}
          </span>
          <div className="remote-users" aria-label="Active collaborators">
            {remoteUsers.length === 0 && (
              <span className="user-pill muted">Just you</span>
            )}
            {remoteUsers.map((user) => (
              <span
                key={user.id}
                className="user-pill"
                style={{ borderColor: user.color, color: user.color }}
                title={user.name || user.id}
              >
                {user.name || user.id}
              </span>
            ))}
          </div>
        </div>
      </div>
      <div ref={editorContainerRef} className="editor-container" />
      <div className="editor-footer">
        Collaborative Markdown editor powered by Yjs & CodeMirror 6.
      </div>
    </div>
  );
};
