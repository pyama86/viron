import ObjectAssign from 'object-assign';
import riot from 'riot';
export default function() {
  const store = this.riotx.get();

  let tag;

  const fadeIn = () => {
    setTimeout(() => {
      this.isVisible = true;
      this.update();
    }, 100);
  };

  const fadeOut = () => {
    this.isVisible = false;
    this.update();
    setTimeout(() => {
      store.action('popovers.remove', this.opts.id);
    }, 1000);
  };

  this.on('mount', () => {
    tag = riot.mount(this.refs.content, this.opts.tagname, ObjectAssign({
      isPopover: true,
      popoverCloser: fadeOut
    }, this.opts.tagopts))[0];
    fadeIn();
    window.addEventListener('keydown', this.handleKeyDown);
    window.addEventListener('resize', this.handleWindowResize);
  }).on('before-unmount', () => {
    tag.unmount(true);
  }).on('unmount', () => {
    window.removeEventListener('keydown', this.handleKeyDown);
    window.removeEventListener('resize', this.handleWindowResize);
  });

  this.handleTap = () => {
    fadeOut();
  };

  this.handleFrameTap = e => {
    // 内部イベントを外部に伝播させない。
    e.stopPropagation();
  };

  this.handleKeyDown = e => {
    switch (e.keyCode) {
    case 27: // Esc
      fadeOut();
      break;
    default:
      break;
    }
  };

  this.handleWindowResize = () => {
    fadeOut();
  };
}
